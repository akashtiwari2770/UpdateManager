# Artillery Monitoring Guide

## How to Check if Artillery is Running

### Method 1: Check Process List (Quick Check)

```bash
# Check if Artillery process exists
ps aux | grep -i artillery | grep -v grep

# Or using pgrep (returns PID if found)
pgrep -f artillery

# More detailed process info
ps aux | grep artillery | grep -v grep
```

**Output examples:**
- If running: Shows process details with PID, CPU, memory usage
- If not running: No output or "No Artillery processes found"

### Method 2: Check Process by PID

If you know the PID (from when you started it):

```bash
# Check if specific PID is running
ps -p <PID>

# Get detailed info about the process
ps -p <PID> -o pid,cmd,etime,stat,pcpu,pmem

# Check process tree
pstree -p <PID>
```

### Method 3: Check Process Status

```bash
# Check if Artillery is running and get status
pgrep -fl artillery

# Check with process details
ps -ef | grep artillery | grep -v grep

# Check process with resource usage
top -p $(pgrep -f artillery | tr '\n' ',' | sed 's/,$//')
```

### Method 4: Monitor Artillery in Real-Time

If Artillery is running in the foreground, you'll see live output. If running in background:

```bash
# Check if running in background
jobs

# Bring background job to foreground
fg %1  # Replace 1 with job number

# Or check the terminal where you started it
```

### Method 5: Check Artillery Output Files

Artillery may create output files:

```bash
# Check for recent Artillery report files
ls -lth /tmp/artillery-* 2>/dev/null

# Check for Artillery logs
find /tmp -name "*artillery*" -mmin -10 2>/dev/null
```

### Method 6: Check Network Activity

Artillery makes HTTP requests, so you can monitor network activity:

```bash
# Monitor network connections (if Artillery is making requests)
netstat -anp | grep artillery

# Or using ss
ss -anp | grep artillery

# Check if target port (8080) has connections
netstat -an | grep :8080 | wc -l
```

### Method 7: Check System Resources

Artillery uses CPU and memory:

```bash
# Check CPU/Memory usage of Artillery
ps aux | grep artillery | grep -v grep | awk '{print $2, $3, $4, $11}'

# Monitor in real-time
watch -n 1 'ps aux | grep artillery | grep -v grep'
```

## Understanding Artillery Execution

### Artillery Test Phases

Artillery tests run in phases. You can tell if it's still running by:

1. **Check the terminal output** - Artillery shows progress:
   ```
   Phase started: Warm up (index: 0, duration: 30s)
   Phase started: Ramp up (index: 1, duration: 60s)
   Phase started: Sustained load (index: 2, duration: 120s)
   ```

2. **Check process status** - If process exists, test is likely running

3. **Monitor network activity** - Active connections indicate test is running

### Artillery Test Durations

Based on your test configurations:

- **Mixed Load Test** (`artillery-config.yml`): ~5 minutes (30+60+120+30+30 seconds)
- **Read-Heavy Test** (`artillery-read-heavy.yml`): ~5 minutes
- **Write-Heavy Test** (`artillery-write-heavy.yml`): ~3 minutes
- **Spike Test** (`artillery-spike-test.yml`): ~1.5 minutes

## Quick Status Check Script

Create a simple script to check Artillery status:

```bash
#!/bin/bash
# Save as: check-artillery.sh

PID=$(pgrep -f artillery)

if [ -z "$PID" ]; then
    echo "❌ Artillery is NOT running"
    exit 1
else
    echo "✅ Artillery IS running (PID: $PID)"
    echo ""
    echo "Process details:"
    ps -p $PID -o pid,cmd,etime,stat,pcpu,pmem
    echo ""
    echo "To stop: kill $PID"
    exit 0
fi
```

Make it executable and use:
```bash
chmod +x check-artillery.sh
./check-artillery.sh
```

## Running Artillery in Background

### Start Artillery in Background

```bash
# Run in background and save output
make load-test > artillery-output.log 2>&1 &

# Or with nohup (survives terminal close)
nohup make load-test > artillery-output.log 2>&1 &

# Get the PID
echo $!
```

### Monitor Background Artillery

```bash
# Check if still running
ps -p <PID>

# Monitor output file
tail -f artillery-output.log

# Check progress
tail -f artillery-output.log | grep "Phase started"
```

## Stopping Artillery

If Artillery is running and you need to stop it:

```bash
# Find the PID
PID=$(pgrep -f artillery)

# Stop gracefully (SIGTERM)
kill $PID

# Or force kill if needed
kill -9 $PID

# Or kill all Artillery processes
pkill -f artillery
```

## Artillery Completion

Artillery completes when:

1. **All phases finish** - You'll see "Summary report" in output
2. **Process exits** - `ps` won't show Artillery process
3. **Final report appears** - Shows test statistics

Example completion output:
```
Summary report @ 14:30:15(+0000) 2025-01-15
  Scenarios launched:  1000
  Scenarios completed: 998
  Requests completed:  5000
  Mean response/sec: 83.33
  ...
```

## Troubleshooting

### Artillery appears stuck

1. **Check if server is responding:**
   ```bash
   curl http://localhost:8080/health
   ```

2. **Check for errors in output:**
   ```bash
   # If running in background, check log file
   tail -f artillery-output.log | grep -i error
   ```

3. **Check system resources:**
   ```bash
   # Check if system is overloaded
   top
   free -h
   ```

4. **Check network connectivity:**
   ```bash
   # Verify target server is reachable
   curl -v http://localhost:8080/health
   ```

### Artillery process exists but no activity

This might indicate:
- Test is between phases (brief pause)
- Network issues (check server)
- System resource constraints

Check with:
```bash
# Monitor network activity
watch -n 1 'netstat -an | grep :8080 | wc -l'
```

## Useful Commands Summary

```bash
# Quick check if running
pgrep -f artillery && echo "Running" || echo "Not running"

# Get PID
pgrep -f artillery

# Check process details
ps -p $(pgrep -f artillery) -o pid,cmd,etime,stat

# Monitor in real-time
watch -n 1 'ps aux | grep artillery | grep -v grep'

# Check network activity
netstat -an | grep :8080

# Stop Artillery
pkill -f artillery

# View recent Artillery output (if logged)
tail -100 artillery-output.log
```

## Integration with Makefile

You can add a check command to your Makefile:

```makefile
# Check if Artillery is running
artillery-status:
	@PID=$$(pgrep -f artillery 2>/dev/null); \
	if [ -z "$$PID" ]; then \
		echo "Artillery is NOT running"; \
	else \
		echo "Artillery IS running (PID: $$PID)"; \
		ps -p $$PID -o pid,cmd,etime,stat; \
	fi

# Stop Artillery
artillery-stop:
	@pkill -f artillery && echo "Artillery stopped" || echo "No Artillery process found"
```

Then use:
```bash
make artillery-status
make artillery-stop
```

