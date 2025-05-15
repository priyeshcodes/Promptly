# 🧾 Promptly Task Manager in Go

A simple yet powerful command-line Task Manager built using Go and the **Cobra CLI** library. It allows you to add, view, update, and delete tasks with support for **priorities**, **deadlines**, and **dependencies**... and many more


## 📌 Features

- Add tasks with title, description, priority, and deadline
- Support for task dependencies
- List tasks with clear formatting
- Update and delete tasks
- Filter by priority or due date
- Data persistence via JSON (or extendable to DB)
- User-friendly terminal commands


## 📸 Screenshot 
![Promptly4](https://github.com/user-attachments/assets/ff500bf6-f355-4c30-9fa2-cc3848f9753b)


### 🔧 Prerequisites

- [Go 1.19+](https://golang.org/dl/)
- Terminal or command line access

### 📥 Installation 
git clone https://github.com/priyeshcodes/Promptly.git
cd Promptly
go run main.go

### ➕ Add a Task
go run main.go add --title "Write blog" --description "Use **Markdown**" --priority high --deadline 2025-05-20

### 🔗 Add Dependencies 
go run main.go add --title "Deploy app" --description "Final step" --priority high --deadline 2025-05-30 --dependencies "Write tests,Optimize DB"

### 📃 List Tasks 
go run main.go list

### 🗑️ Delete Task
go run main.go delete --title "Update Docs"

### Task By Id.
go run main.go complete f64dd42c-46ad-4e98-a801-abc123... //just a id here it can be different

### 📊 Task Statistics.
go run main.go stats

### 🔔 Notify
go run main.go notify

### Edit Task
go run main.go edit 123abc //this is just an id 123abc

### 🔍 Search
go run main.go search golang // basically a Grep



### Finally Promptly
go run main.go tui //just press <Enter> and see the magic itself

### Coming Soon!
# A full enhanced Frontend
# Colored Output
# Task Tagging

### 🪪 License
This project is open-source under the MIT License.

