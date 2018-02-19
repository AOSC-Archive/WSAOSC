#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <shellapi.h>

#include <stdio.h>
#include <wchar.h>
#include <stdbool.h>
#include <ctype.h>

#include "wslutil.h"

#define DISTOR_NAME L"AOSC"

void print_help(void)
{
	puts(
		"WSAOSC Utility\n"
		"Usage: aosc command\n"
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
			// run without args
			break;
		}

		wchar_t *command = argv[1];
		if (wcscmp(command, L"run") == 0 || wcscmp(command, L"exec") == 0)
		{
		}
		else if (wcscmp(command, L"config") == 0)
		{
		}
		else if (wcscmp(command, L"uninstall") == 0 || wcscmp(command, L"clean") == 0)
		{
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
