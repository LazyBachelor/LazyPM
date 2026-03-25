# Project Management Interface Survey - Participant Guide

Welcome! This survey evaluates how people use different project management interfaces to complete common tasks.

## Overview

**Duration**: 15-20 minutes
**Goal**: Complete 8 task-based exercises using different interfaces and provide feedback on your experience

## Data Collection Notice

You will be given a username and password you have to enter before starting the survey.
We do this to ensure no unwanted actors are interfering with our data collecting.

- Responses are anonymized before analysis
- We collect task results, questionnaire answers, and interaction events during the survey
- We do not intentionally collect personally identifying information
- You may exit at any time by pressing Esc between tasks or closing the terminal
- Data is collected continuously, so exiting early still contributes valuable insights
- Contact us if you want your data deleted.

## After Each Task

You'll answer a short questionnaire about:
- How difficult the interface was to use
- How difficult the task itself was
- Specific questions about the task experience

## Tips for Success

1. **Task Details**: At the start of each task, a "Task Details.txt" file is created with full instructions.
   You can have this open throughout the survey and go back to it any time you need.
2. **Take Your Time**: There's no time limit, focus on completing the task correctly
3. **Read Carefully**: Each task has specific requirements that will be validated
4. **Exit Anytime**: Press Esc or q (depending on interface) to exit early if needed
5. **Interface Help**: Each interface has built-in help - look for it!

## Troubleshooting

- **Shift Key Stuck**: There is currently a bug where the shift key gets stuck,
   if this happens it usually fixes itself after tapping shift.
- **No database connection**: You can continue without submitting responses and submit later with `pm submit`
- **Terminal too small**: Resize your terminal window to at least 55x16 characters,
   we recommend to keeping the terminal as large as possible for best experience.
- **Interface not responding**: Press Ctrl+C to exit, then restart the survey



## Interfaces You'll Use

You will complete tasks using different interfaces. Here's what to expect:

### REPL (Read-Eval-Print Loop)
- Interactive command environment
- Type `pm help` to see available commands
- To see task progress use the command `status` or `pm status`
- Suggestions appear as you type use tab to auto complete
- You can run shell commands directly (e.g., `git`, `ls`, `cat`)
- Type `exit` to skip task during the survey

### TUI (Terminal User Interface)
- Text-based interactive interface
- Navigate using keyboard inputs
- For task progress press `?` to view details.
- Help menu at bottom shows available keybinds
- Two views: List View (browse issues) and Kanban (track progress)
- Press `q` to quit task during the survey

### Web Interface
- Accessed through a web browser
- Mouse-driven interaction only
- For task progress press the red button in the upper right corner.
- Two views: Issues (list format) and Kanban (sprint tracking)
- Press `esc` or `q` in the terminal to skip the task



## The Tasks

You will complete 8 tasks in sequential order with different interfaces:

### 1. Create Issue Task
**Difficulty**: Medium | **Time**: ~3 minutes

Create a new issue with specific details:
1. Create an issue titled: **"My first Issue"**
2. Add description: **"I need to do some coding"**
3. Assign to: **"Me"**
4. Mark as: **"In Progress"**



### 2. Backlog Refinement Task
**Difficulty**: Easy | **Time**: ~2 minutes

Clean up duplicate issues in the backlog:
1. Go to the backlog
2. Find two issues with the same name
3. Open one of them
4. Select **"Close issue"**
5. Choose **"Duplicate issue"** as closing reason
6. Save/close the issue



### 3. Sprint Planning Task
**Difficulty**: Medium | **Time**: ~3 minutes

Organize the backlog for a new sprint:
1. Review backlog issues
2. Select issues to include in the sprint
3. Add issues to the sprint
4. Address blocked/dependent issues
5. Prioritize high-priority items

**Goal**: Create a realistic sprint plan with at least 3 issues, including at least 2 high-priority items.



### 4. Priority Management Task
**Difficulty**: Medium | **Time**: ~4 minutes

Rebalance sprint priorities due to a critical production issue:

A database connection issue needs urgent attention:
1. Find the database-related issue
2. Change its priority to **4** (highest)
3. Set all other issues (features and chores) to priority **1**



### 5. Dependency Management Task
**Difficulty**: Medium | **Time**: ~3 minutes

Manage issue dependencies:
1. Find 2 issues that mention dependencies in their description
2. Set their status to **"Blocked"**
3. Find the foundational issue that others depend on
4. Set its priority to **3**
5. Set its status to **"In Progress"**
6. Assign it to: **"Me"**



### 6. Issue Review and Cleanup Task
**Difficulty**: Medium | **Time**: ~3 minutes

Review and maintain project issues:
1. Add a comment to **two** different issues
2. Delete the issue titled **"Delete this issue"**



### 7. Coding Task
**Difficulty**: Hard | **Time**: ~5 minutes

Fix a logical error while managing an issue:

1. A text file (**code.txt**) will appear in your directory
2. Open it and follow the instructions inside
3. Fix the code: change the function to add two numbers (`a + b`) instead of subtract (`a - b`)
4. Save the file
5. Change the status of the issue **"Fix major error"** to **"Closed"**

### 8. Git Task
**Difficulty**: Hard | **Time**: ~5 minutes

Perform a Git operation:
1. Assign the issue **"Upgrade the codebase"** to **"Me"**
2. A folder called **"task"** is created in your directory
3. Inside, find **README.md**
4. Edit the file (add your name or a short note)
5. Commit your changes:
   - Open a terminal in the task folder
   - Run: `git add .`
   - Run: `git commit -m "updated codebase"`
6. Set the issue status to **"Closed"**


## Thank You!

Your participation helps us understand how different interfaces impact project management workflows.
We appreciate your time and feedback!
