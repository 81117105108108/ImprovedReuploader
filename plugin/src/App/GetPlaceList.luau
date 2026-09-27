--!strict
-- Shared place-list accessor. Filter tab registers its getter; asset tabs
-- consume it. Replaces _G.getPlaceList (global + race if Filter hasn't run yet).
local getter: (() -> { number })? = nil

local GetPlaceList = {}

function GetPlaceList.set(fn: () -> { number })
	getter = fn
end

function GetPlaceList.get(): { number }
	while getter == nil do
		task.wait()
	end
	return (getter :: () -> { number })()
end

return GetPlaceList
