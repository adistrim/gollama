# Project Gollama

This chat app is built on top of the [go-openai](https://github.com/sashabaranov/go-openai) library.  
It talks to the local **Ollama** server by default, but you can also hook it up to **OpenAI** or **DeepSeek**.  
Basically: if it works with the OpenAI SDK, it works here.

It’s not just a pretty chat box, the bot can also connect to a GitHub (using MCP tools). That means it can:  
- Peek at your issues (even private ones)  
- Work on them  
- And open a pull request from a fresh branch 

I’ve mostly tested this with [`gpt-oss:20b`](https://ollama.com/library/gpt-oss), and it’s been surprisingly good at tool use.  
Still early days though, so don’t expect magic on mid or large codebases just yet.

It needs a GitHub access token to work.

---

## Works with
These all I have tried
- ✅ Local Ollama models (Qwen, LLaMA, gpt-oss)  
- ✅ OpenAI  
- ✅ DeepSeek  
- ✅ Basically anything OpenAI-compatible (this i haven't) 

---

## Setup guide

### 1. pull the code
```bash
git clone https://github.com/adistrim/gollama
cd gollama
````

### 2. Backend (Go)

```bash
cd server
go mod tidy             # install dependencies
cp .env.example .env    # put your stuff here
go build -o gollama .   # build the server
./gollama               # run the server
```

### 3. Frontend (React)

```bash
cd ../app
pnpm install            # install node dependencies
cp .env.example .env
pnpm dev                # dev server
```

Other frontend scripts:

```bash
pnpm dev      # dev server
pnpm build    # prod build
pnpm preview  # preview build
```

If you get stuck, open an issue or just email me [araj@adistrim.in](mailto:araj@adistrim.in)

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
