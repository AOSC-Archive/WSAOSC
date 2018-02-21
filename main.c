#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <shellapi.h>
#include <objbase.h>

#include <stdio.h>
#include <wchar.h>
#include <stdbool.h>
#include <ctype.h>
#include <stdlib.h>

#ifdef _MSC_VER
#define _CRT_SECURE_NO_WARNINGS
#pragma warning(disable:4996)
#endif /* _MSC_VER */

#include "wslutil.h"
#include "my_gets.h"

#define DISTOR_NAME L"AOSC-OS"

void print_help(void)
{
	puts(
		"AOSC OS for WSL Utility 0.1.1\n"
		"Usage: aosc-os [<command>]\n\n"
		"Commands:\n"
		"  run/exec [<command line>] - Run command line in WSL.\n\n"
		"  config [<option> <value>] - Configure or show distor settings.\n"
		"    options:\n"
		"      enable-interop <true/false>\n"
		"      append-nt-path <true/false>\n"
		"      enable-drive-mounting <true/false>\n"
		"      default-uid <uid>\n"
		"      default-user <username> - set only\n\n"
		"  uninstall/clean - Uninstall the distro.\n\n"
		"  help - Print this message."
	);
}

HRESULT new_user(void)
{
	HRESULT hr;
	DWORD exit_code;
	char username[33];
	wchar_t command[64];
	do {
		fputs("Enter new UNIX username: ", stdout);
		my_gets_s(username, sizeof(username));

		if (strlen(username) == 0)
			continue;

		swprintf(command, ARRAYSIZE(command), L"useradd -m %hs", username);

		hr = _WslLaunchInteractive(DISTOR_NAME, command, FALSE, &exit_code);
	} while (FAILED(hr) || exit_code != 0);

	swprintf(command, ARRAYSIZE(command), L"usermod -aG sudo %hs", username);
	hr = _WslLaunchInteractive(DISTOR_NAME, command, FALSE, &exit_code);
	if (SUCCEEDED(hr) && exit_code == 0)
	{
		swprintf(command, ARRAYSIZE(command), L"passwd %hs", username);
		do {
			hr = _WslLaunchInteractive(DISTOR_NAME, command, FALSE, &exit_code);
		} while (FAILED(hr) || exit_code != 0);

		wchar_t username_w[33];
		for (int i = 0; i < ARRAYSIZE(username_w) && username[i]; i++)
			username_w[i] = username[i];

		ULONG uid = uid_from_username(DISTOR_NAME, username_w);
		if (uid != ULONG_MAX)
			hr = _WslConfigureDistribution(DISTOR_NAME, uid, WSL_DISTRIBUTION_FLAGS_DEFAULT);
	}
	return hr;
}

bool install_distor(void)
{
	fputs("Will install AOSC OS for WSL, continue? [y/N] ", stdout);

	char ch = getchar();
	char c;
	while ((c = getchar()) != '\n' && c != EOF);
	if (tolower(ch) == 'y')
	{
		puts("Installing, this may take a few minutes...");

		HRESULT hr = _WslRegisterDistribution(DISTOR_NAME, L"install.tar.gz");
		if (SUCCEEDED(hr))
		{
			DWORD exit_code;
			hr = _WslLaunchInteractive(DISTOR_NAME, L"groupadd sudo && sed -i 's/^#\\s*\\(%sudo\\)/\\1/' /etc/sudoers", FALSE, &exit_code);
			if (SUCCEEDED(hr))
			{
				hr = new_user();
			}
		}

		if (SUCCEEDED(hr))
		{
			puts("Installation successful!");
			return true;
		}
		else
		{
			printf("Installation failed! (0x%X)\n", hr);
		}
	}
	else
	{
		puts("Abort.");
	}
	return false;
}

bool uninstall_distor(void)
{
	fputs("Will uninstall AOSC OS for WSL, continue? [y/N] ", stdout);
	char ch = getchar();
	if (tolower(ch) == 'y')
	{
		puts("Removing filesystem...");
		HRESULT hr = _WslUnregisterDistribution(DISTOR_NAME);
		if (SUCCEEDED(hr))
		{
			puts("Successfully removed distro.");
			return true;
		}
		else
		{
			printf("Failed to remove distor! (0x%X)\n", hr);
		}
	}
	else
	{
		puts("Abort.");
	}
	return false;
}

int run_wsl(int argc, wchar_t *argv[])
{
	size_t command_len = 0;
	for (int i = 0; i < argc; i++)
		command_len += wcslen(argv[i]) + 1;

	wchar_t *command_line = NULL;
	if (command_len != 0)
	{
		command_line = (wchar_t*)malloc(command_len * sizeof(wchar_t));
		if (!command_line)
			return false;
		command_line[0] = L'\0';

		for (int i = 0; i < argc; i++)
		{
			wcscat(command_line, argv[i]);
			size_t len = wcslen(command_line);

			if (i < (argc - 1))
			{
				command_line[len] = L' ';
				command_line[len + 1] = L'\0';
			}
		}
		command_line[command_len - 1] = L'\0';
	}

	DWORD exit_code;
	HRESULT hr = _WslLaunchInteractive(DISTOR_NAME, command_line, TRUE, &exit_code);

	free(command_line);

	if (FAILED(hr))
	{
		printf("Error: 0x%X\n", hr);
		return 1;
	}

	return exit_code;
}

bool config(int argc, wchar_t *argv[])
{
	if (argc < 2)
	{
		ULONG default_uid;
		WSL_DISTRIBUTION_FLAGS distor_flags;
		if (!get_wsl_distro_config(DISTOR_NAME, &default_uid, &distor_flags))
			return false;

		bool enable_interop = (distor_flags & WSL_DISTRIBUTION_FLAGS_ENABLE_INTEROP);
		bool append_nt_path = (distor_flags & WSL_DISTRIBUTION_FLAGS_APPEND_NT_PATH);
		bool enable_drive_mounting = (distor_flags & WSL_DISTRIBUTION_FLAGS_ENABLE_DRIVE_MOUNTING);

		printf(
			"enable-interop %s\n"
			"append-nt-path %s\n"
			"enable-drive-mounting %s\n"
			"default-uid %lu\n",
			enable_interop ? "true" : "false",
			append_nt_path ? "true" : "false",
			enable_drive_mounting ? "true" : "false",
			default_uid);
	}
	else
	{
		wchar_t *option = argv[0];
		wchar_t *value = argv[1];

		WSL_DISTRIBUTION_FLAGS flag = 0;
		ULONG uid;
		if (wcscmp(option, L"enable-interop") == 0)
			flag = WSL_DISTRIBUTION_FLAGS_ENABLE_INTEROP;
		else if (wcscmp(option, L"append-nt-path") == 0)
			flag = WSL_DISTRIBUTION_FLAGS_APPEND_NT_PATH;
		else if (wcscmp(option, L"enable-drive-mounting") == 0)
			flag = WSL_DISTRIBUTION_FLAGS_ENABLE_DRIVE_MOUNTING;
		else if (wcscmp(option, L"default-uid") == 0)
		{
			uid = wcstoul(value, NULL, 10);
		}
		else if (wcscmp(option, L"default-user") == 0)
		{
			uid = uid_from_username(DISTOR_NAME, value);
			if (uid == ULONG_MAX)
			{
				printf("Can not get UID from username '%ls'!", value);
				return false;
			}
		}
		else
		{
			printf("Invaild option '%ls'!", option);
			return false;
		}

		bool enable = false;
		if (flag)
		{
			if (wcscmp(value, L"true") == 0)
				enable = true;
			else if (wcscmp(value, L"false") == 0)
			{
				// pass
			}
			else
			{
				printf("Invaild value '%ls' for option '%ls'!", value, option);
				return false;
			}
		}

		ULONG default_uid;
		WSL_DISTRIBUTION_FLAGS distor_flags;
		if (!get_wsl_distro_config(DISTOR_NAME, &default_uid, &distor_flags))
			return false;

		if (flag)
		{
			if (enable)
				distor_flags |= flag;
			else
				distor_flags &= ~flag;
		}
		else
			default_uid = uid;

		HRESULT hr = _WslConfigureDistribution(DISTOR_NAME, default_uid, distor_flags);

		return SUCCEEDED(hr);
	}
	return true;
}

int main(void)
{
	if (!load_wsl_api())
	{
		printf("Failed to load wsl api! (%lu)\n", GetLastError());
		return 1;
	}

	int argc;
	LPWSTR *argv = CommandLineToArgvW(GetCommandLineW(), &argc);
	if (argv == NULL)
	{
		printf("CommandLineToArgvW failed! (%lu)\n", GetLastError());
		return 1;
	}

	int retval = 0;
	do {
		if (!_WslIsDistributionRegistered(DISTOR_NAME))
		{
			if (!install_distor())
			{
				retval = 1;
				break;
			}
		}

		if (argc < 2)
		{
			retval = run_wsl(0, NULL);
			break;
		}

		wchar_t *command = argv[1];
		if (wcscmp(command, L"run") == 0 || wcscmp(command, L"exec") == 0)
		{
			retval = run_wsl(argc - 2, argv + 2);
		}
		else if (wcscmp(command, L"config") == 0)
		{
			if (!config(argc - 2, argv + 2))
				retval = 1;
		}
		else if (wcscmp(command, L"uninstall") == 0 || wcscmp(command, L"clean") == 0)
		{
			if (!uninstall_distor())
				retval = 1;
		}
		else
		{
			print_help();
			retval = 1;
		}
	} while (0);

	LocalFree(argv);

	return retval;
}
