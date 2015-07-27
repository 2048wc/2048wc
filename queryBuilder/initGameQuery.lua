redis.call("FLUSHALL")
redis.call("RPUSH", KEYS[1], "prevGameid");
redis.call("RPUSH", "prevGameid", "finish");
local lastGameid = redis.call("LRANGE", KEYS[1], -1, -1); 
if lastGameid[1] == nil then
	error("Failed to find the last gameid of this user"); 
	return false; 
end 
local lastGameMove = redis.call("LRANGE", lastGameid[1], -1, -1); 
if lastGameMove[1] == nil then 
	error("Failed to find the last move of the last game of this user") 
	return false; 
end 
if lastGameMove[1] == "finish" then 
   	local call1 = redis.call("RPUSH", KEYS[1], KEYS[3]); 
   	local call2 = redis.call("RPUSH", KEYS[3], KEYS[2]);
    local checkuserid = redis.call("LRANGE", KEYS[1], -1, -1);
    local checkgameMove = redis.call("LRANGE", KEYS[3], -1, -1);
    print("call1", call1);
    print("call2", call2);
    print("checkuserid", checkuserid[1]);
    print("checkGameMove", checkgameMove[1]); 
	return true; 
else 
    return false; 
end
