(window.webpackJsonp=window.webpackJsonp||[]).push([[13],{506:function(e,a,s){"use strict";s.r(a);var t=s(45),n=Object(t.a)({},(function(){var e=this,a=e.$createElement,s=e._self._c||a;return s("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[s("h1",{attrs:{id:"installing-3rd-party-libraries"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#installing-3rd-party-libraries"}},[e._v("#")]),e._v(" Installing 3rd party libraries "),s("Badge",{attrs:{text:"experimental",type:"warning"}})],1),e._v(" "),s("p",[e._v("The ABS interpreter comes with a built-in installer for 3rd party libraries,\nvery similar to "),s("code",[e._v("npm install")]),e._v(", "),s("code",[e._v("pip install")]),e._v(" or "),s("code",[e._v("go get")]),e._v(".")]),e._v(" "),s("p",[e._v("The installer, bundled since the "),s("code",[e._v("1.8.0")]),e._v(" release, is currently "),s("strong",[e._v("experimental")]),e._v("\nand a few things might change.")]),e._v(" "),s("p",[e._v("In order to install a package, you simply need to run "),s("code",[e._v("abs get")]),e._v(":")]),e._v(" "),s("div",{staticClass:"language-bash extra-class"},[s("pre",{pre:!0,attrs:{class:"language-bash"}},[s("code",[e._v("$ abs get github.com/abs-lang/abs-sample-module \n🌘  - Downloading archive\nUnpacking"),s("span",{pre:!0,attrs:{class:"token punctuation"}},[e._v("..")]),e._v(".\nCreating alias"),s("span",{pre:!0,attrs:{class:"token punctuation"}},[e._v("..")]),e._v(".\nInstall Success. You can use the module with "),s("span",{pre:!0,attrs:{class:"token variable"}},[s("span",{pre:!0,attrs:{class:"token variable"}},[e._v("`")]),e._v("require"),s("span",{pre:!0,attrs:{class:"token punctuation"}},[e._v("(")]),s("span",{pre:!0,attrs:{class:"token string"}},[e._v('"abs-sample-module"')]),s("span",{pre:!0,attrs:{class:"token punctuation"}},[e._v(")")]),s("span",{pre:!0,attrs:{class:"token variable"}},[e._v("`")])]),e._v("\n")])])]),s("p",[e._v("Modules will be saved under the "),s("code",[e._v("vendor/$MODULE")]),e._v(" directory. Each module\nalso gets an alias to facilitate requiring them in your code, meaning that\nboth of these forms are supported:")]),e._v(" "),s("div",{staticClass:"language- extra-class"},[s("pre",{pre:!0,attrs:{class:"language-text"}},[s("code",[e._v('⧐  require("abs-sample-module/sample.abs")\n{"another": f() {return hello world;}}\n\n⧐  require("vendor/github.com/abs-lang/abs-sample-module/sample.abs")\n{"another": f() {return hello world;}}\n')])])]),s("p",[e._v("Module aliases are saved in the "),s("code",[e._v("packages.abs.json")]),e._v(" file\nwhich is created in the same directory where you run the\n"),s("code",[e._v("abs get ...")]),e._v(" command:")]),e._v(" "),s("div",{staticClass:"language- extra-class"},[s("pre",{pre:!0,attrs:{class:"language-text"}},[s("code",[e._v('$ abs get github.com/abs-lang/abs-sample-module\n🌗  - Downloading archive\nUnpacking...\nCreating alias...\nInstall Success. You can use the module with `require("abs-sample-module")`\n\n$ cat packages.abs.json \n{\n    "abs-sample-module": "./vendor/github.com/abs-lang/abs-sample-module"\n}\n')])])]),s("p",[e._v("If an alias is already taken, the installer will let you know that you\nwill need to use the full path when requiring the module:")]),e._v(" "),s("div",{staticClass:"language- extra-class"},[s("pre",{pre:!0,attrs:{class:"language-text"}},[s("code",[e._v('$ echo \'{"abs-sample-module": "xyz"}\' > packages.abs.json \n\n$ abs get github.com/abs-lang/abs-sample-module\n🌘  - Downloading archive\nUnpacking...\nCreating alias...This module could not be aliased because module of same name exists\n\nInstall Success. You can use the module with `require("./vendor/github.com/abs-lang/abs-sample-module")`\n')])])]),s("p",[e._v("When requiring a module, ABS will try to load the "),s("code",[e._v("index.abs")]),e._v(" file unless\nanother file is specified:")]),e._v(" "),s("div",{staticClass:"language- extra-class"},[s("pre",{pre:!0,attrs:{class:"language-text"}},[s("code",[e._v('$ ~/projects/abs/builds/abs                                          \nHello alex, welcome to the ABS programming language!\nType \'quit\' when you\'re done, \'help\' if you get lost!\n\n⧐  require("abs-sample-module")\n{"another": f() {return hello world;}}\n\n⧐  require("abs-sample-module/index.abs")\n{"another": f() {return hello world;}}\n\n⧐  require("abs-sample-module/another.abs")\nf() {return hello world;}\n')])])]),s("h2",{attrs:{id:"supported-hosting-platforms"}},[s("a",{staticClass:"header-anchor",attrs:{href:"#supported-hosting-platforms"}},[e._v("#")]),e._v(" Supported hosting platforms")]),e._v(" "),s("p",[e._v("Currently, the installer supports modules hosted on:")]),e._v(" "),s("ul",[s("li",[e._v("GitHub")])])])}),[],!1,null,null,null);a.default=n.exports}}]);