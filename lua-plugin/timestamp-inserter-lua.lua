local core = require("apisix.core")
local http = require("resty.http")

local plugin_name = "timestamp-inserter-lua"

local schema = {
    type = "object",
    properties = {
        timestamp_service_uri = {
            type = "string",
            description = "URL of the timestamp service"
        }
    },
    required = {"timestamp_service_uri"}
}

local _M = {
    version = 0.1,
    priority = 4002,
    name = plugin_name,
    schema = schema
}

function _M.check_schema(conf)
    local ok, err = core.schema.check(schema, conf)
    if not ok then
        return false, err
    end
    return true
end

local function fetchTimestamp(conf)
    local httpc = http.new()

    -- Call the timestamp service
    local res, err = httpc:request_uri(conf.timestamp_service_uri, {
        method = "GET"
    })

    if not res then
        core.log.error("Failed to call timestamp service: ", err)
        return nil, "Failed to fetch timestamp"
    end

    if res.status ~= 200 then
        core.log.error("Failed to fetch timestamp: ", res.status, res.body)
        return nil, "Error fetching timestamp"
    end

    local timestamp = res.body
    return timestamp, nil
end

function _M.rewrite(conf, ctx)
    -- Fetch the timestamp from the timestamp service
    local timestamp, err = fetchTimestamp(conf)
    if not timestamp then
        core.log.error("Error fetching timestamp: ", err)
        return
    end

    -- received timestamp looks like this = {"time":"2025-04-18T10:42:25+05:30"}
    -- I want to extract the time value and then set it in the request header
    local time_value = core.json.decode(timestamp)
    if not time_value or not time_value.time then
        core.log.error("Invalid timestamp format: ", timestamp)
        return
    end
    -- Extract the time value
    local time = time_value.time
    -- Log the timestamp
    core.log.info("Fetched timestamp: ", time)
    -- Set the timestamp in the request header
    core.request.set_header(ctx, "X-Timestamp", time)
end

return _M