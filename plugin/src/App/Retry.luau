--!strict
-- Variadic generics (A..., R...) need very new Luau; keep simple ...any for parser compat.
return function(tries: number, func: (...any) -> ...any, ...: any): (boolean, ...any)
	local result: any
	for i=1, tries do
		result = { pcall(func, ...) }
		local success = result[1]
		if success then break end
		task.wait(0.5 * i)
	end
	return result[1], table.unpack(result, 2)
end
