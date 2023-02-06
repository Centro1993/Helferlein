# Helferlein

This tool uses ChatGPT to Describe, Create and Execute Shell Commands.
```
Usage: helferlein <arg> "<prompt>"
-c, c, create:
        Describe what you want to achieve, get a Linux Command in return
-d, d, describe:
        Enter a Linux command, get a Description of what it does
-h, h, help:
        Display this Help Section
```
Before you run Helferlein, set your API-Token by running "export OPENAI_API_KEY=<TOKEN>"
Create a Token at https://platform.openai.com/account/api-keys

## Usage Example
```
helferlein c "Read current CPU Usage"
top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1"%"}'

1) Get a Description
2) Run the Command
3) Exit
WARNING: Don't blindly run commands you do not understand.

1 / 2 / 3: 1

This command prints the percentage of CPU usage that is not being used.
Do you want to run this command? (y/n)
y

Output:
0%
```

## WARNING
This is a stupid Idea, and you should not be using it. Don't blame ChatGPT, OpenAI or me if you wreck your Stuff ;)
