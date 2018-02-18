#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <shellapi.h>

#include <stdio.h>
#include <wchar.h>
#include <stdbool.h>
#include <ctype.h>

#include "wslutil.h"

void print_help(void)
{
	puts(
		"WSAOSC Utility\n"
		"Usage: aosc command\n"
	);
}

void install_distor(void)
{
	fputs("Will install AOSC for WSL, continue? [y/N]", stdout);
	char ch = getchar();
	if (tolower(ch) == 'y')
	{

	}
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

	do {
		if (argc < 2)
		{
			print_help();
			break;
		}

		wchar_t *command = argv[1];
		if (wcscmp(command, L"run") == 0 || wcscmp(command, L"exec") == 0)
		{
		}
		else if (wcscmp(command, L"config") == 0)
		{
		}
		else if (wcscmp(command, L"install") == 0)
		{
			install_distor();
		}
		else if (wcscmp(command, L"uninstall") == 0 || wcscmp(command, L"clean") == 0)
		{
		}
		else
		{
			print_help();
		}

	} while (0);

	LocalFree(argv);

	return 0;
}
