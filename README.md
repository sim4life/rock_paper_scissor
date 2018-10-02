# Rock Paper Scissor  
Work on goroutines, channels &amp; concurrency

Write a multithreaded program that represents two players playing _Rock, Paper, Scissors_.  
Each player is modeled as an independently executing thread.  
There is a main thread which controls the game and gets the next move from each player thread.  
The main thread produces a JSON file that represents the game history.  
The final output should be a JSON file that contains the complete game history of all game rounds since the start of the program:

```json
 [
   {
      "Round": 1,
      "Winner": "Player2",
      "Inputs": {
        "Player1": "rock",
        "Player2": "paper"
      }
   }, {
      "Round": 2,
      "Winner": null,
      "Inputs": {
        "Player1": "scissors",
        "Player2": "scissors"
      }
    }, ... {
      "Round": 100,
      "Winner": "Player1",
      "Inputs": {
        "Player1": "scissors",
        "Player2": "paper"
      }
    }, 
]
```

## Components
The program must have the following:  
1. A thread for Player1 with the following properties  
  i. It has the ability to receive a request from another thread (in this case, the main thread) to return it’s next move.  
  ii. The input it receives should include the full history of the previous rounds with Player2, which should look like the above example  
  iii. It can employ any strategy to generate it’s move, including choosing a random move. It’s up to you decide what strategy to use.  
2. A thread for Player2 with the following properties . 
  i. It’s identical to Player2, except it employs a different strategy for generating it’s move. Player2’s strategy should be based on the history of Player1’s previous moves. For example, randomly choosing one of Player1’s moves and using that as it’s next move. Or, using Player1’s last move as it’s next move.  
3. A main thread that does the following:  
  i. Start the Player1 and Player2 threads  
  ii. Repeat the following steps 100 times:    
    a. Ask Player1 and Player2 for their next move each, passing in the full history of previous rounds (initially empty)  
    b. Determine the winner for this round (if any)  
    c. Store the results in memory, which will be passed to the players as inputs on subsequent rounds  
  iii. Emit the results in a file called result.json and output a message to the console with the full path to the result file.  
  iv. Stop the Player1 and Player2 threads  

## Testing
Write a test that exercises the code and validates that the JSON output is valid and contains 100 game rounds
