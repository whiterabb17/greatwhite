# GreatWhite / Necromancy

[![CodeFactor](https://www.codefactor.io/repository/github/whiterabb17/greatwhite/badge)](https://www.codefactor.io/repository/github/whiterabb17/greatwhite)

### Synopsis
GreatWhite / Necromancy
Some musing into developing a simple GUI for a multi-os reverse trojan, built entirely (mostly) using Golang.
Building a TeamServer and Client in Golang is rudimentry. But I figured there had to be a way, after having done alot of standalone apps in electron just before I wanted to make a Golang frontend that wasnt as unnecessary complex as something like Wails. Which don't get me wrong, its an amazing choice for Golang frontends! But the amount of extra generated code seemed needless to me. I wanted to make it easy to add and remove functions in either the template html or the golang handling code without needing to entirely regenerate all the code every time, to create all the new binding-fields in the go-code.

A natural choice for golang GUIs is vue or react, which again is easy with fair stack-dev knowledge.
You can also use things like mux or gin which are another 2 obvious options for my endevour. Though, i've done my share of full-stack and my reason for this project was learning new things. 

Which i did in extremely copious amounts ..... x_x

So after learning <a href="https://github.com/fyne-io/fyne">Fyne</a>. Which is a great framework, allows for a container-like build scheme, making what i see akin to something like the subzero trojan GUI(?) possibly, but with a much more sleek and futuristic look. Though not the right choice for this project - it was the GUI framework i used to create a red-team WireShark-like MITM packet-sniffer which does allow for dynamic changes to the forms real-time.

Then there was <a href="https://github.com/gioui/gio">Gio</a>, which again i think was created very well! Allowing it to become quite easy over time by just remembering the templating pattern - to start building beautiful looking GUIs fairly quickly.

After having built and rebuilt my GUI over and over and over - which i think aged me alot in the process. 
I eventually settled on <a href="https://github.com/asticode/go-astilectron">go-astilectron</a>.
Which is essentially a Golang port of electron. While i suppose, because its an electron port it doesnt exactly full under `entirely Golang`
I felt it was a close enough compromise - as it allowed me to take advantage of some of my favourite electron-modules, namely <a href="https://github.com/sindresorhus/electron-context-menu">electron-context-menu</a>, giving me dynamic WinForm-like right-click context menus which was an amazing + for something essentially running in a chrome-sandbox
Other then that some other pros where that:
- it uses html templates (which was an original want)
- it has an event-handler system (like electron obviously) which i was already familiar with
- the rest of the rendering and backend functionality could be coding in Golang and any part could be edited without the need to re-bind functions from the front to the back.
- re-binding is only necessary in the cause of static resource changes AND only when in the final packaged application - which if done right only needs to be packaged at release.
- as well as hot-reload development to the html templates which is even more convienient because you can do your html/css dev real-time and quickly make backend function changed without sitting for another 5mins after an edit to re-bind a function.

While go-astilectron is very powerful in your ability to freely create each of your components in a way that is easy to pick up with some knowledge into other fairly standard frontend languages (html / css / javascript) which is actually a pretty big advantage with all other frameworks i mentioned above.<br>
It does have a initialization system that i'll admit took me longer then i'd have liked to admit to get a grip on.<br>
There is, for this purpose a convienient <a href="https://github.com/asticode/go-astilectron-bundler">bundler</a> and <a href="https://github.com/asticode/go-astilectron-bootstrap">bootstrapper</a> that will essentially handle all the packaging of static files and resources and general templating which is extremely handly if you just need a single window application!<br>
As soon as i wanted extra windows.. it became a slightly more troublesome thing. Which admittedly is why i then ended up learning the Gio framework.<br>
Then after a light-bulb moment - i realized my problem which helped me get enough of a grip to create my final package without the bundler.<br>
Though myself and the developer both encourage you to use the bundler and bootstrapper and it does make things alot easier!

    i think sometimes i just like making myself suffer. Im sure there was a word for that? Nevertheless.

My main goal was mostly creating an (almost) fully Golang frontend i needed something to generate the data to fill it.<br>
So because i had been making some backdoors and red-team tools and was already in GoDev mode, i made the go-client, which is just a minimal-backdoor and its chromium recovery modules in Go. Then made extra agents in CSharp and Rust because i wanted to get a better hang of multi-language functionality and code interoperability.
Maybe i'll release the code for those clients at some-point.
<br><br>
### NOTICE
I am releasing this publicly for entirely education purposes.<br>
Due to that reason i have purposefully bugged certain parts of the GUI as to stop skiddies from abusing this code.<br>
As well as defanged the client and removed the recovery modules, again to stop skids from recklessness.<br>
<br>
<img src="greatwhite.gif" />
