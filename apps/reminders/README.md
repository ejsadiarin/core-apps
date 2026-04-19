# Reminder tool for various topics (for spaced repetition)

## Topics

- `ps` tool with `awk` practice

```bash
ps -eo pid,ppid,stat,cmd | awk '$4 == "sleep" {print $2}' | xargs -I {} kill -9 {}
ps aux --sort=%mem | tail -n 6
```

- CKA exam questions (see DumbITGuy playlist, killerkoda, etc.)

- daily quests as todos (workout, programming)
- anime/manga/manhwa/manhua updates
- tech/dev daily news
- backlog clearing/old reads to revisit
- "Core reminds me to..." (philosophical, grounded mental)
- self learn, have strong fundamentals (language-agnostic)
- job hop for growth potential
- the work you do needs to be broadly communicated for visibility in order to be considered for high ratings and compensation increases.
- programming is just moving data around and managing state.
- academic/organization work
- job application/openings (opt, since i have linkedin, etc.)

### for reference

- budget tracker application
  - record expenses everyday
  - also record gains (from allowance, salary, revenue, assets etc.)
- reminder tool
  - daily quests as todos (workout, programming)
  - anime/manga/manhwa/manhua updates
  - tech/dev daily news
  - backlog clearing/old reads to revisit
  - "Core reminds me to..." (philosophical, grounded mental)
    - self learn, have strong fundamentals (language-agnostic)
    - job hop for growth potential
    - the work you do needs to be broadly communicated for visibility in order to be considered for high ratings and compensation increases.
    - programming is just moving data around and managing state.
  - academic/organization work
  - job application/openings (opt, since i have linkedin, etc.)
- wardrobe management system
  - start picturing clothes with name (like in uniqlo), description
- learning management system (quiz or LMS)
  - catered to myself.
  - Course -> Lesson -> Topic -> Question -> AnswerOption
  - see: [High-level approach: two buckets - content vs user state](https://www.reddit.com/r/Backend/comments/1qlz9h6/comment/o1i9yns/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button)
