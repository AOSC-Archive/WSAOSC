#ifndef _WSL_UTIL
#define _WSL_UTIL

/* Determines if a distribution is already registered */
typedef BOOL (WINAPI *fnWslIsDistributionRegistered)(PCWSTR distributionName);
fnWslIsDistributionRegistered _WslIsDistributionRegistered;

/* Registers a new distribution given the information provided. */
typedef HRESULT (WINAPI *fnWslRegisterDistribution)(PCWSTR distributionName, PCWSTR tarGzFilename);
fnWslRegisterDistribution _WslRegisterDistribution;

/* Unregisters the specified distribution */
typedef HRESULT (WINAPI *fnWslUnregisterDistribution)(PCWSTR distributionName);
fnWslUnregisterDistribution _WslUnregisterDistribution;

/* Flags specifying WSL behavior */
typedef enum
{
	WSL_DISTRIBUTION_FLAGS_NONE = 0x0,
	WSL_DISTRIBUTION_FLAGS_ENABLE_INTEROP = 0x1,
	WSL_DISTRIBUTION_FLAGS_APPEND_NT_PATH = 0x2,
	WSL_DISTRIBUTION_FLAGS_ENABLE_DRIVE_MOUNTING = 0x4
} WSL_DISTRIBUTION_FLAGS;

#define WSL_DISTRIBUTION_FLAGS_VALID (WSL_DISTRIBUTION_FLAGS_ENABLE_INTEROP | WSL_DISTRIBUTION_FLAGS_APPEND_NT_PATH | WSL_DISTRIBUTION_FLAGS_ENABLE_DRIVE_MOUNTING)
#define WSL_DISTRIBUTION_FLAGS_DEFAULT (WSL_DISTRIBUTION_FLAGS_ENABLE_INTEROP | WSL_DISTRIBUTION_FLAGS_APPEND_NT_PATH | WSL_DISTRIBUTION_FLAGS_ENABLE_DRIVE_MOUNTING)

/* Configure the given distribution */
typedef HRESULT (WINAPI *fnWslConfigureDistribution)(PCWSTR distributionName, ULONG defaultUID, WSL_DISTRIBUTION_FLAGS wslDistributionFlags);
fnWslConfigureDistribution _WslConfigureDistribution;

/* Get the given distribution's configuration info */
typedef HRESULT (WINAPI *fnWslGetDistributionConfiguration)(PCWSTR distributionName, ULONG * distributionVersion, ULONG * defaultUID, WSL_DISTRIBUTION_FLAGS * wslDistributionFlags, PSTR ** defaultEnvironmentVariables, ULONG * defaultEnvironmentVariableCount);
fnWslGetDistributionConfiguration _WslGetDistributionConfiguration;

typedef HRESULT (WINAPI *fnWslLaunchInteractive)(PCWSTR distributionName, PCWSTR command, BOOL useCurrentWorkingDirectory, DWORD * exitCode);
fnWslLaunchInteractive _WslLaunchInteractive;

typedef HRESULT (WINAPI *fnWslLaunch)(PCWSTR distributionName, PCWSTR command, BOOL useCurrentWorkingDirectory, HANDLE stdIn, HANDLE stdOut, HANDLE stdErr, HANDLE * process);
fnWslLaunch _WslLaunch;

bool load_wsl_api(void)
{
	HMODULE wsl_handle = LoadLibraryExW(L"wslapi.dll", NULL, LOAD_LIBRARY_SEARCH_SYSTEM32);
	if (!wsl_handle)
		return false;

	_WslIsDistributionRegistered = (fnWslIsDistributionRegistered)GetProcAddress(wsl_handle, "WslIsDistributionRegistered");
	_WslRegisterDistribution = (fnWslRegisterDistribution)GetProcAddress(wsl_handle, "WslRegisterDistribution");
	_WslUnregisterDistribution = (fnWslUnregisterDistribution)GetProcAddress(wsl_handle, "WslUnregisterDistribution");
	_WslConfigureDistribution = (fnWslConfigureDistribution)GetProcAddress(wsl_handle, "WslConfigureDistribution");
	_WslGetDistributionConfiguration = (fnWslGetDistributionConfiguration)GetProcAddress(wsl_handle, "WslGetDistributionConfiguration");
	_WslLaunchInteractive = (fnWslLaunchInteractive)GetProcAddress(wsl_handle, "WslLaunchInteractive");
	_WslLaunch = (fnWslLaunch)GetProcAddress(wsl_handle, "WslLaunch");

	if (!_WslIsDistributionRegistered ||
		!_WslRegisterDistribution ||
		!_WslUnregisterDistribution ||
		!_WslConfigureDistribution ||
		!_WslGetDistributionConfiguration ||
		!_WslLaunchInteractive ||
		!_WslLaunch)
		return false;

	return true;
}

bool get_wsl_distro_config(const wchar_t *distor_name, ULONG *default_uid, WSL_DISTRIBUTION_FLAGS *distor_flags)
{
	ULONG distor_version, _default_uid, env_count;
	WSL_DISTRIBUTION_FLAGS _distor_flags;
	PSTR *env;
	HRESULT hr = _WslGetDistributionConfiguration(distor_name, &distor_version, &_default_uid, &_distor_flags, &env, &env_count);
	if (FAILED(hr))
		return false;

	for (ULONG i = 0; i < env_count; i++)
		CoTaskMemFree(env[i]);
	CoTaskMemFree(env);

	if (default_uid)
		*default_uid = _default_uid;

	if (distor_flags)
		*distor_flags = _distor_flags;

	return true;
}

ULONG uid_from_username(const wchar_t *distor_name, const wchar_t *username)
{
	HANDLE read_pipe, write_pipe;
	SECURITY_ATTRIBUTES sec_attr;
	sec_attr.nLength = sizeof(SECURITY_ATTRIBUTES);
	sec_attr.lpSecurityDescriptor = NULL;
	sec_attr.bInheritHandle = TRUE;

	if (CreatePipe(&read_pipe, &write_pipe, &sec_attr, 0))
	{
		wchar_t command[64];
		swprintf(command, ARRAYSIZE(command), L"id -u %s", username);

		HANDLE handle_stdin = GetStdHandle(STD_INPUT_HANDLE);
		HANDLE handle_stderr = GetStdHandle(STD_ERROR_HANDLE);
		HANDLE process_handle;
		HRESULT hr = _WslLaunch(distor_name, command, FALSE, handle_stdin, write_pipe, handle_stderr, &process_handle);
		if (SUCCEEDED(hr))
		{
			DWORD exit_code;
			WaitForSingleObject(process_handle, INFINITE);
			GetExitCodeProcess(process_handle, &exit_code);
			CloseHandle(process_handle);

			if (exit_code == 0)
			{
				char uid_string[16];
				DWORD byte_read;
				if (ReadFile(read_pipe, uid_string, 15, &byte_read, NULL))
				{
					CloseHandle(read_pipe);
					CloseHandle(write_pipe);

					uid_string[byte_read] = '\0';
					ULONG uid = strtoul(uid_string, NULL, 10);

					return uid;
				}
			}
		}
	}
	return ULONG_MAX;
}

#endif /* _WSL_UTIL */
