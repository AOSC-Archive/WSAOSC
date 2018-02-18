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

#endif /* _WSL_UTIL */
