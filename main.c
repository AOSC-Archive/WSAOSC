#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <shellapi.h>

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

#define DISTOR_NAME L"AOSC"

void print_help(void)
{
	puts(
		"WSAOSC Utility\n"
		"Usage: aosc [<command>]\n\n"
		"Commands:\n"
		"  run/exec [<command line>] - Run command line in WSL.\n\n"
		"  config [<option> <value>] - Configure or show distor settings.\n"
		"    options:\n"
		"      enable-interop <true/false>\n"
		"      append-nt-path <true/false>\n"
		"      enable-drive-mounting <true/false>\n\n"
		"  uninstall/clean - Uninstall the distro.\n\n"
		"  help - Print this message."
	);
}

bool install_distor(void)
{
	fputs("Will install AOSC for WSL, continue? [y/N] ", stdout);

	char ch = getchar();
	if (tolower(ch) == 'y')
	{
		puts("Installing, this may take a few minutes...");

		HRESULT hr = _WslRegisterDistribution(DISTOR_NAME, L"install.tar.gz");

		if (SUCCEEDED(hr))
		{
			hr = _WslConfigureDistribution(DISTOR_NAME, 0, WSL_DISTRIBUTION_FLAGS_DEFAULT);
		}

		if (SUCCEEDED(hr))
		{
			puts("Installation successful!");
			return true;
		}
		else
		{
			printf("Installation failed! (%X)\n", hr);
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
	fputs("Will uninstall AOSC for WSL, continue? [y/N] ", stdout);
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
			printf("Failed to remove! (%X)\n", hr);
		}
	}
	else
	{
		puts("Abort.");
	}
	return false;
}

bool run_wsl(int argc, wchar_t *argv[])
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
	return SUCCEEDED(hr);
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
			if (!run_wsl(0, NULL))
				retval = 1;
			break;
		}

		wchar_t *command = argv[1];
		if (wcscmp(command, L"run") == 0 || wcscmp(command, L"exec") == 0)
		{
			if (!run_wsl(argc - 2, argv + 2))
				retval = 1;
		}
		else if (wcscmp(command, L"config") == 0)
		{
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
