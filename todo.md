Add systems(delete too)
Bind systems.


Adding works, bindings can be added but they have 0 functionality so far

# 🌌 Ecosystem & League System Development Summary

## ✅ Current Priority Order

1. **Ecosystem Manager (in progress)**
   - Manages systems, bindings, and execution
   - In final stretch (working on binding removal UI)

2. **League Training System** (next)
   - Goal: structure and log deliberate League practice
   - Will connect to ecosystem manager and emit events

3. **Calendar / Todo System** (after League)
   - Organizes days, training cycles, system reminders
   - Receives events from other systems

4. **Kayaking System**
   - Later-stage project involving API data (e.g., Garmin)
   - Creative, exploratory — reward after core systems

---

## 🛠️ League System: Design Plan

### 🎯 Goal
Support structured, skill-focused training like athletic practice.

### 📦 Core Components
- **Focus Areas** (e.g. disruption, ganking, resource mgmt)
- **Training Logs** (manually logged sessions with focus, games played, notes)
- **Training Plan** (define cycles like 3 training + 1 test day)
- **Review Mode** (see last X days: what worked, common issues)
- **Optional Meta Tracking** (energy, focus, tilt, etc.)

### 🧰 CLI Example Commands
```bash
eco run league log --focus disruption --games 3 --notes "..."
eco run league review --days 2
eco run league plan --cycle 7 --focus-order disruption weakpoints ganking
eco run league today

