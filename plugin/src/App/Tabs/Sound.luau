--!strict
local getFilterOptions = require("../GetFilterOptions")
local reuploadIds = require("../ReuploadIds")
local UiLibrary = require("../../UiLibrary")
local GetPlaceList = require("../GetPlaceList")

local Selection = game:GetService("Selection")

return function(ui: UiLibrary.Ui, plugin: Plugin)
    local tab = ui:CreateTab("Sound")

    tab:CreateButton("Reupload", function()
        local filter = getFilterOptions(plugin, game:GetDescendants())
        table.insert(filter.WhitelistedInstances, "Sound")

		reuploadIds(plugin, ui, filter, "Sound", GetPlaceList.get()) -- forgive me
	end)
	
	tab:CreateButton("Reupload Selected", function()
        local filter = {
            WhitelistedInstances = { "Sound", "StringValue", "NumberValue", "IntValue", "LuaSourceContainer" },
            Instances = Selection:Get()
        }

        reuploadIds(plugin, ui, filter, "Sound", GetPlaceList.get()) -- forgive me
	end)


end
