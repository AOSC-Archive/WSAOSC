HydroGo = require("HydroGo")
sh = require("sh")
Hydro = require("Hydro")

function config()
    print(HydroGo.detectGlide())
    print(HydroGo.DetectGlideYaml())
    HydroGo.setShowAll(true)
    -- HydroGo.setBuildAll(true)
    if HydroGo.DetectGlideYaml() == false then
        HydroGo.InitDeps()
        HydroGo.InstallDeps()
    end
end

function HydroDev.default()
	target.build()
end

function target.build()
    Env.Path()
    config()
    HydroGo.build()
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
