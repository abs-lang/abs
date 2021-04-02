<template>
  <div class="wrapper">
    <div class="message">
      <h3>ABS Playground</h3>
      <p>
        This is a playground for you to test ABS code directly in the browser:
        you can execute code by either pressing <code>Ctrl+Enter</code> or
        clicking on the button right below the editor.
      </p>
      <div class="editor-wrapper">
        <MonacoEditor
          class="editor"
          v-model="code"
          language="shell"
          :options="options"
          @editorDidMount="editorDidMount"
          ref="editor"
        />
        <div class="results">
          <h4>Output</h4>
          <pre class="output-shell" v-text="out"></pre>
          <h4>Result</h4>
          <pre class="output-shell" v-text="result"></pre>
        </div>
        <div class="button-wrapper">
          <button class="run-button" @click="run">RUN</button>
        </div>
      </div>
      <p>
        Please note that <a href="/syntax/system-commands/">system commands</a> like
        <code>ls -la</code> do not work as this code is running directly in your
        browser and not on a server, but you can still explore ABS’ syntax
        without having to download and install it in your system.
      </p>
    </div>
  </div>
</template>

<script>
import MonacoEditor from "vue-monaco";

let mod, inst, go;

async function fetch_wasm_module() {
  const response = await fetch("/abs.wasm");
  const buffer = await response.arrayBuffer();
  const obj = await WebAssembly.instantiate(buffer, go.importObject);
  mod = obj.module;
  inst = obj.instance;
  await go.run(inst);
}

export default {
  components: {
    MonacoEditor,
  },
  beforeCreate() {
    let s = document.createElement("script");
    s.setAttribute("type", "application/javascript");
    s.setAttribute("src", "/wasm_exec.js");
    s.onload = function load() {
      go = new Go();
      fetch_wasm_module();
    };
    document.body.appendChild(s);

  },
  methods: {
    editorDidMount() {
      let run_code = this.run
      this.$refs.editor.getEditor().addAction({
        // An unique identifier of the contributed action.
        id: 'my-unique-id',

        // A label of the action that will be presented to the user.
        label: 'My Label!!!',

        // An optional array of keybindings for the action.
        keybindings: [
          monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter,
        ],

        // A precondition for this action.
        precondition: null,

        // A rule to evaluate on top of the precondition in order to dispatch the keybindings.
        keybindingContext: null,

        contextMenuGroupId: 'navigation',

        contextMenuOrder: 1.5,

        // Method that will be executed when the action is triggered.
        // @param editor The editor instance is passed in as a convinience
        run: function(ed) {
          run_code()
        }
      });
    },
    run() {
      if (!go) {
        go = new Go();
        fetch_wasm_module();
      }
      let { out, result } = abs_run_code(this.code);
      this.out = out;
      this.result = result;
    },
  },

  data() {
    return {
      out: `Here is where everything that's outputted from the script will appear. Try running "echo(123)"!`,
      result: `Here is where the return value of the script will appear. Try running "return 42"!`,
      code: `# Try pressing ctrl + enter
lebron = {
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

return lebron.id
      `,
      options: {
        fontSize: 12.5,
        fontFamily: "Source Code Pro",
        minimap: {
          enabled: false,
        },
      },
    };
  },
};
</script>

<style scoped>
.red {
  color: #ee6680;
}
.wrapper {
  display: grid;
  grid-template-columns: 1fr;
}
.editor-wrapper {
  margin-top: 40px;
  border-bottom: 2px solid #ccc;
  border: 2px solid #eee;
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  grid-template-rows: 1fr 50px;
  box-shadow: 0 0 10px #eee inset;
  padding: 3px;
}
.button-wrapper {
  border-top: 2px solid #eee;
  padding: 15px;
  grid-row: 2;
  grid-column: 1/3;
  justify-content: center;
  align-items: center;
}
.run-button {
  background-color: #ee6680;
  color: white;
  border: 1px solid #ee6680;
  padding: 5px 15px;
  cursor: pointer;
  font-weight: bold;
}
.output-shell {
  padding: 20px;
  background: #333;
  color: rgb(253 255 0 / 80%);
  border-radius: 5px;
  overflow-y: scroll;
  margin-right: 10px;
  font-size: 12px;
  height: 120px;
  overflow-x: hidden;
  overflow-y: hidden;
}
.results {
  display: grid;
  grid-template-rows: 20px 1fr 20px 1fr;
  row-gap: 5px;
}
.results h4 {
  margin: 0;
}

pre {
    white-space: pre-wrap;       /* Since CSS 2.1 */
    white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
    white-space: -pre-wrap;      /* Opera 4-6 */
    white-space: -o-pre-wrap;    /* Opera 7 */
    word-wrap: break-word;       /* Internet Explorer 5.5+ */
}
</style>