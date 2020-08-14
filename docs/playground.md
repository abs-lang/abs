# ABS playground

This is a playground for you to test ABS code directly
in the browser: you can execute code by either pressing
`Ctrl+Enter` or clicking on the button right below the
editor.

Please note that [system commands](/syntax/system-commands)
like `` `ls -la` `` do not work as this code is running directly in your
browser and not on a server, but you can still explore ABS' syntax
without having to download and install it in your system.

<script src="/wasm_exec.js"></script>

<script>
  const go = new Go();
  let mod, inst;

  async function fetch_wasm_module() {
    const response = await fetch("/abs.wasm");
    const buffer = await response.arrayBuffer();
    const obj = await WebAssembly.instantiate(buffer, go.importObject);
    mod = obj.module;
    inst = obj.instance;
    await go.run(inst);
  }
  fetch_wasm_module()

  function run() {
    var editor = ace.edit("editor");
    let code = editor.getValue();
    let {out, result} = abs_run_code(code)
    document.getElementById("out-area").innerHTML = out || "No output"
    document.getElementById("result-area").innerHTML = result || "No return value"
  }
</script>

<style type="text/css" media="screen">
    #editor { 
        width: 100%;
        height: 350px;
    }
    .ace_content{
      left: 20px;
      padding-top: 20px;
    }
    .ace_layer {
      top: 20px;
    }
</style>

<div id="editor">lebron = {
  "id": 23, 
  "name": "LeBron James", 
  "nicknames": [
    "the king", "king james", "the chosen one"
  ]
}

echo("%s is also known as:", lebron.name)

for nickname in lebron.nicknames {
  echo("* %s", nickname.title())
}

return lebron.id</div>

<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.4.5/ace.js" type="text/javascript" charset="utf-8"></script>
<script>
    var editor = ace.edit("editor");
    editor.setTheme("ace/theme/monokai");
    editor.session.setMode("ace/mode/javascript");
    editor.getSession().setUseWorker(false);
    editor.setOptions({
      fontSize: "13pt"
    });
    editor.commands.addCommand({
        name: "run",
        bindKey: {win: "Ctrl-Enter", mac: "Command-Enter"},
        exec: function(editor) {
          run()
        }
    });
</script>

<button onClick="run();" id="runButton">Run code with Ctrl+Enter, or click HERE</button>

## Output

<pre style="white-space: pre-line" id="out-area">
Here is where everything that's outputted from the script will appear. Try running "echo(123)"!
</pre>

## Result

<pre style="white-space: pre-line" id="result-area">
Here is where the return value of the script will appear. Try running "return 42"!
</pre>

## Next

That's about it for this section!

You can now head over to read about ABS's syntax,
starting with [assignments](/syntax/assignments)!