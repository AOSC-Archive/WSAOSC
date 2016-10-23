HydroGo = require("HydroGo")
sh = require("sh")
Hydro = require("Hydro")

function config()
    print(HydroGo.detectGlide())
    print(HydroGo.DetectGlideYaml())
    HydroGo.setShowAll(true)
    HydroGo.setTargetOS("windows")
    HydroGo.setTargetArch("amd64")
    -- HydroGo.setBuildAll(true)
    if HydroGo.DetectGlideYaml() == false then
        HydroGo.InitDeps()
        HydroGo.InstallDeps()
    else
        HydroGo.InstallDeps()
    end
    if Common.detectInPath("rsrc") == false then
        print(color("%{magenta}The rsrc tool not found in path. The produced executable may not run!%{reset}\n"))
    else 
        print(color("%{green}The rsrc is found in path.%{reset}\n"))
    end
    HydroDev.sh("rsrc -manifest WSAOSC.exe.manifest -ico aosc.ico WSAOSC.syso")
end

function HydroDev.default()
	target.build()
end

function target.build()
    Env.Path()
    config()
    HydroGo.build()
end

function target.run()
    HydroDev.sh("./WSAOSC.exe")
end

function target.clean()
    HydroDev.sh("rm -f glide.*")
    HydroDev.sh("*.exe")
    HydroDev.sh("rm -rf vendor")
end

function target.rebuild()
    target.clean()
    HydroDev.default()
end

function target.watch()
	HydroDev.watch{callback=target.build, dir=".", recursive=true, filter=".go", exclude={".exe"}}
end
