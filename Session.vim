let SessionLoad = 1
if &cp | set nocp | endif
let s:cpo_save=&cpo
set cpo&vim
inoremap <silent> <Plug>(neocomplete_auto_refresh) =neocomplete#mappings#refresh()
inoremap <silent> <Plug>(neocomplete_start_manual_complete) =neocomplete#mappings#manual_complete()
inoremap <silent> <Plug>(neocomplete_start_auto_complete) =neocomplete#mappings#auto_complete()
inoremap <silent> <Plug>(neocomplete_start_omni_complete) 
inoremap <silent> <expr> <Plug>(neocomplete_start_unite_quick_match) unite#sources#neocomplete#start_quick_match()
inoremap <silent> <expr> <Plug>(neocomplete_start_unite_complete) unite#sources#neocomplete#start_complete()
nnoremap <silent>  :TmuxNavigateLeft
nnoremap <silent> <NL> :TmuxNavigateDown
nnoremap <silent>  :TmuxNavigateUp
nnoremap <silent>  :TmuxNavigateRight
map  <Plug>(ctrlp)
nmap  :w
nnoremap <silent>  :TmuxNavigatePrevious
nmap ,i <Plug>(go-info)
nmap ,s <Plug>(go-implements)
nmap ,gv <Plug>(go-doc-vertical)
nmap ,gd <Plug>(go-doc)
nmap ,dt <Plug>(go-def-tab)
nmap ,dv <Plug>(go-def-vertical)
nmap ,ds <Plug>(go-def-split)
nmap ,c <Plug>(go-coverage)
nmap ,t <Plug>(go-test)
nmap ,b <Plug>(go-build)
nmap ,r <Plug>(go-run)
nnoremap ,, 
nnoremap K :Ag ''
xmap S <Plug>VSurround
nnoremap <silent> [b :bprevious
nnoremap <silent> ]b :bnext
nmap cS <Plug>CSurround
nmap cs <Plug>Csurround
nmap ds <Plug>Dsurround
vmap gx <Plug>NetrwBrowseXVis
nmap gx <Plug>NetrwBrowseX
xmap gS <Plug>VgSurround
nmap ySS <Plug>YSsurround
nmap ySs <Plug>YSsurround
nmap yss <Plug>Yssurround
nmap yS <Plug>YSurround
nmap ys <Plug>Ysurround
nnoremap <silent> <Plug>(go-pointsto) :call go#guru#PointsTo(-1)
nnoremap <silent> <Plug>(go-alternate-split) :call go#alternate#Switch(0, "split")
nnoremap <silent> <Plug>(go-alternate-vertical) :call go#alternate#Switch(0, "vsplit")
nnoremap <silent> <Plug>(go-alternate-edit) :call go#alternate#Switch(0, "edit")
nnoremap <silent> <Plug>(go-vet) :call go#lint#Vet(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-lint) :call go#lint#Golint()
nnoremap <silent> <Plug>(go-metalinter) :call go#lint#Gometa(0)
nnoremap <silent> <Plug>(go-doc-browser) :call go#doc#OpenBrowser()
nnoremap <silent> <Plug>(go-doc-split) :call go#doc#Open("new", "split")
nnoremap <silent> <Plug>(go-doc-vertical) :call go#doc#Open("vnew", "vsplit")
nnoremap <silent> <Plug>(go-doc-tab) :call go#doc#Open("tabnew", "tabe")
nnoremap <silent> <Plug>(go-doc) :call go#doc#Open("new", "split")
nnoremap <silent> <Plug>(go-def-stack-clear) :call go#def#StackClear()
nnoremap <silent> <Plug>(go-def-stack) :call go#def#Stack()
nnoremap <silent> <Plug>(go-def-pop) :call go#def#StackPop()
nnoremap <silent> <Plug>(go-def-tab) :call go#def#Jump("tab")
nnoremap <silent> <Plug>(go-def-split) :call go#def#Jump("split")
nnoremap <silent> <Plug>(go-def-vertical) :call go#def#Jump("vsplit")
nnoremap <silent> <Plug>(go-def) :call go#def#Jump('')
nnoremap <silent> <Plug>(go-rename) :call go#rename#Rename(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-sameids-toggle) :call go#guru#ToggleSameIds()
nnoremap <silent> <Plug>(go-whicherrs) :call go#guru#Whicherrs(-1)
nnoremap <silent> <Plug>(go-sameids) :call go#guru#SameIds()
nnoremap <silent> <Plug>(go-referrers) :call go#guru#Referrers(-1)
nnoremap <silent> <Plug>(go-channelpeers) :call go#guru#ChannelPeers(-1)
xnoremap <silent> <Plug>(go-freevars) :call go#guru#Freevars(0)
nnoremap <silent> <Plug>(go-callstack) :call go#guru#Callstack(-1)
nnoremap <silent> <Plug>(go-describe) :call go#guru#Describe(-1)
nnoremap <silent> <Plug>(go-callers) :call go#guru#Callers(-1)
nnoremap <silent> <Plug>(go-callees) :call go#guru#Callees(-1)
nnoremap <silent> <Plug>(go-implements) :call go#guru#Implements(-1)
nnoremap <silent> <Plug>(go-imports) :call go#fmt#Format(1)
nnoremap <silent> <Plug>(go-import) :call go#import#SwitchImport(1, '', expand('<cword>'), '')
nnoremap <silent> <Plug>(go-info) :call go#tool#Info(0)
nnoremap <silent> <Plug>(go-deps) :call go#tool#Deps()
nnoremap <silent> <Plug>(go-files) :call go#tool#Files()
nnoremap <silent> <Plug>(go-coverage-browser) :call go#coverage#Browser(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-coverage-toggle) :call go#coverage#BufferToggle(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-coverage-clear) :call go#coverage#Clear()
nnoremap <silent> <Plug>(go-coverage) :call go#coverage#Buffer(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-test-compile) :call go#test#Test(!g:go_jump_to_error, 1)
nnoremap <silent> <Plug>(go-test-func) :call go#test#Func(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-test) :call go#test#Test(!g:go_jump_to_error, 0)
nnoremap <silent> <Plug>(go-install) :call go#cmd#Install(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-generate) :call go#cmd#Generate(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-build) :call go#cmd#Build(!g:go_jump_to_error)
nnoremap <silent> <Plug>(go-run) :call go#cmd#Run(!g:go_jump_to_error)
vnoremap <silent> <Plug>NetrwBrowseXVis :call netrw#BrowseXVis()
nnoremap <silent> <Plug>NetrwBrowseX :call netrw#BrowseX(expand((exists("g:netrw_gx")? g:netrw_gx : '<cfile>')),netrw#CheckIfRemote())
nnoremap <silent> <Plug>SurroundRepeat .
nnoremap <silent> <Plug>(ctrlp) :CtrlP
nnoremap <SNR>19_: :=v:count ? v:count : ''
imap S <Plug>ISurround
imap s <Plug>Isurround
inoremap <expr> 	 pumvisible() ? "\" : "\	"
imap  <Plug>Isurround
let &cpo=s:cpo_save
unlet s:cpo_save
set autoindent
set backspace=indent,eol,start
set backupdir=~/.vim/backup//
set cmdheight=3
set completeopt=preview,menuone
set confirm
set directory=~/.vim/swap//
set errorformat=%-G#\ %.%#,%-G%.%#panic:\ %m,%Ecan't\ load\ package:\ %m,%A%f:%l:%c:\ %m,%A%f:%l:\ %m,%C%*\\s%m,%-G%.%#
set expandtab
set fileencodings=ucs-bom,utf-8,default,latin1
set grepprg=ag\ --nogroup\ --nocolor
set helplang=en
set hidden
set history=10000
set hlsearch
set ignorecase
set incsearch
set laststatus=2
set mouse=a
set pastetoggle=<F11>
set printoptions=paper:letter
set ruler
set runtimepath=~/.vim,~/.vim/bundle/Vundle.vim,~/.vim/bundle/vim-fugitive,~/.vim/bundle/vim-go,~/.vim/bundle/vim-rails,~/.vim/bundle/vim-colorschemes,~/.vim/bundle/neocomplete.vim,~/.vim/bundle/gocode/vim/,~/.vim/bundle/ctrlp.vim,~/.vim/bundle/ctrlp-cmatcher,~/.vim/bundle/ag.vim,~/.vim/bundle/vim-ruby,~/.vim/bundle/vim-ansible-yaml,~/.vim/bundle/vim-surround,~/.vim/bundle/vim-gutentags,~/.vim/bundle/editorconfig-vim,~/.vim/bundle/vim-tmux-navigator,/var/lib/vim/addons,/usr/share/vim/vimfiles,/usr/share/vim/vim80,/usr/share/vim/vimfiles/after,/var/lib/vim/addons/after,~/.vim/after,~/.vim/bundle/Vundle.vim,~/.vim/bundle/Vundle.vim/after,~/.vim/bundle/vim-fugitive/after,~/.vim/bundle/vim-go/after,~/.vim/bundle/vim-rails/after,~/.vim/bundle/vim-colorschemes/after,~/.vim/bundle/neocomplete.vim/after,~/.vim/bundle/gocode/vim//after,~/.vim/bundle/ctrlp.vim/after,~/.vim/bundle/ctrlp-cmatcher/after,~/.vim/bundle/ag.vim/after,~/.vim/bundle/vim-ruby/after,~/.vim/bundle/vim-ansible-yaml/after,~/.vim/bundle/vim-surround/after,~/.vim/bundle/vim-gutentags/after,~/.vim/bundle/editorconfig-vim/after,~/.vim/bundle/vim-tmux-navigator/after
set shell=/bin/bash
set shiftwidth=2
set smartcase
set nostartofline
set suffixes=.bak,~,.swp,.o,.info,.aux,.log,.dvi,.bbl,.blg,.brf,.cb,.ind,.idx,.ilg,.inx,.out,.toc
set noswapfile
set tabstop=2
set undodir=~/.vim/undo//
set wildmenu
set window=63
let s:so_save = &so | let s:siso_save = &siso | set so=0 siso=0
let v:this_session=expand("<sfile>:p")
silent only
cd ~/go/src/github.com/ilikeorangutans/phts
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
set shortmess=aoO
badd +1 ~/go/src/github.com/ilikeorangutans/phts
badd +133 db/migrate/1_initial_schema.up.sql
badd +10 db/migrate/2_rendition_configs_on_shares.up.sql
badd +25 db/share_record.go
badd +7 db/record.go
badd +78 db/share_db.go
badd +63 db/share_rendition_configuration_db.go
badd +10 db/album_record.go
badd +40 db/db.go
badd +14 db/paginator.go
badd +27 db/rendition_record.go
badd +8 model/share.go
badd +24 model/share_repository.go
badd +1 main.go
badd +0 admin/api/collections.go
badd +7 test/integration/model/photo_repository_test.go
badd +0 test/integration/model/share_repository_test.go
argglobal
silent! argdel *
argadd ~/go/src/github.com/ilikeorangutans/phts
edit test/integration/model/share_repository_test.go
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd _ | wincmd |
split
wincmd _ | wincmd |
split
2wincmd k
wincmd w
wincmd w
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe '1resize ' . ((&lines * 20 + 32) / 64)
exe 'vert 1resize ' . ((&columns * 159 + 159) / 319)
exe '2resize ' . ((&lines * 19 + 32) / 64)
exe 'vert 2resize ' . ((&columns * 159 + 159) / 319)
exe '3resize ' . ((&lines * 19 + 32) / 64)
exe 'vert 3resize ' . ((&columns * 159 + 159) / 319)
exe 'vert 4resize ' . ((&columns * 159 + 159) / 319)
argglobal
edit /usr/share/vim/vim80/doc/starting.txt
setlocal keymap=
setlocal noarabic
setlocal autoindent
setlocal backupcopy=
setlocal balloonexpr=
setlocal nobinary
setlocal nobreakindent
setlocal breakindentopt=
setlocal bufhidden=
setlocal nobuflisted
setlocal buftype=help
setlocal nocindent
setlocal cinkeys=0{,0},0),:,0#,!^F,o,O,e
setlocal cinoptions=
setlocal cinwords=if,else,while,do,for,switch
setlocal colorcolumn=
setlocal comments=s1:/*,mb:*,ex:*/,://,b:#,:%,:XCOMM,n:>,fb:-
setlocal commentstring=/*%s*/
setlocal complete=.,w,b,u,t,i
setlocal concealcursor=nc
setlocal conceallevel=2
setlocal completefunc=
setlocal nocopyindent
setlocal cryptmethod=
setlocal nocursorbind
setlocal nocursorcolumn
set cursorline
setlocal cursorline
setlocal define=
setlocal dictionary=
setlocal nodiff
setlocal equalprg=
setlocal errorformat=
setlocal expandtab
if &filetype != 'help'
setlocal filetype=help
endif
setlocal fixendofline
setlocal foldcolumn=0
setlocal nofoldenable
setlocal foldexpr=0
setlocal foldignore=#
setlocal foldlevel=0
setlocal foldmarker={{{,}}}
setlocal foldmethod=manual
setlocal foldminlines=1
setlocal foldnestmax=20
setlocal foldtext=foldtext()
setlocal formatexpr=
setlocal formatoptions=tcroql
setlocal formatlistpat=^\\s*\\d\\+[\\]:.)}\\t\ ]\\s*
setlocal formatprg=
setlocal grepprg=
setlocal iminsert=2
setlocal imsearch=2
setlocal include=
setlocal includeexpr=
setlocal indentexpr=
setlocal indentkeys=0{,0},:,0#,!^F,o,O,e
setlocal noinfercase
setlocal iskeyword=!-~,^*,^|,^\",192-255
setlocal keywordprg=
setlocal nolinebreak
setlocal nolisp
setlocal lispwords=
setlocal nolist
setlocal makeprg=
setlocal matchpairs=(:),{:},[:]
setlocal modeline
setlocal nomodifiable
setlocal nrformats=bin,octal,hex
set number
setlocal nonumber
setlocal numberwidth=4
setlocal omnifunc=
setlocal path=
setlocal nopreserveindent
setlocal nopreviewwindow
setlocal quoteescape=\\
setlocal readonly
set relativenumber
setlocal norelativenumber
setlocal norightleft
setlocal rightleftcmd=search
setlocal noscrollbind
setlocal shiftwidth=2
setlocal noshortname
setlocal signcolumn=auto
setlocal nosmartindent
setlocal softtabstop=0
setlocal nospell
setlocal spellcapcheck=[.?!]\\_[\\])'\"\	\ ]\\+
setlocal spellfile=
setlocal spelllang=en
setlocal statusline=
setlocal suffixesadd=
setlocal noswapfile
setlocal synmaxcol=3000
if &syntax != 'help'
setlocal syntax=help
endif
setlocal tabstop=8
setlocal tagcase=
setlocal tags=
setlocal textwidth=78
setlocal thesaurus=
setlocal noundofile
setlocal undolevels=-123456
setlocal nowinfixheight
setlocal nowinfixwidth
setlocal wrap
setlocal wrapmargin=0
silent! normal! zE
let s:l = 1324 - ((19 * winheight(0) + 10) / 20)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
1324
normal! 059|
lcd ~/go/src/github.com/ilikeorangutans/phts
wincmd w
argglobal
nnoremap <buffer> <silent>  :call go#def#StackPop(v:count1)
nnoremap <buffer> <silent> ] :call go#def#Jump("split")
nnoremap <buffer> <silent>  :call go#def#Jump("split")
nnoremap <buffer> <silent>  :GoDef
nnoremap <buffer> <silent> K :GoDoc
xnoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('v', 'prev')
onoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('o', 'prev')
nnoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('n', 'prev')
xnoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('v', 'next')
onoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('o', 'next')
nnoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('n', 'next')
xnoremap <buffer> <silent> af :call go#textobj#Function('a')
onoremap <buffer> <silent> af :call go#textobj#Function('a')
let s:cpo_save=&cpo
set cpo&vim
nnoremap <buffer> <silent> g<LeftMouse> <LeftMouse>:GoDef
nnoremap <buffer> <silent> gd :GoDef
xnoremap <buffer> <silent> if :call go#textobj#Function('i')
onoremap <buffer> <silent> if :call go#textobj#Function('i')
nnoremap <buffer> <silent> <C-LeftMouse> <LeftMouse>:GoDef
let &cpo=s:cpo_save
unlet s:cpo_save
setlocal keymap=
setlocal noarabic
setlocal autoindent
setlocal backupcopy=
setlocal balloonexpr=
setlocal nobinary
setlocal nobreakindent
setlocal breakindentopt=
setlocal bufhidden=
setlocal buflisted
setlocal buftype=
setlocal nocindent
setlocal cinkeys=0{,0},0),:,0#,!^F,o,O,e
setlocal cinoptions=
setlocal cinwords=if,else,while,do,for,switch
setlocal colorcolumn=
setlocal comments=s1:/*,mb:*,ex:*/,://
setlocal commentstring=//\ %s
setlocal complete=.,w,b,u,t,i
setlocal concealcursor=
setlocal conceallevel=0
setlocal completefunc=
setlocal nocopyindent
setlocal cryptmethod=
setlocal nocursorbind
setlocal nocursorcolumn
set cursorline
setlocal cursorline
setlocal define=
setlocal dictionary=
setlocal nodiff
setlocal equalprg=
setlocal errorformat=%-G#\ %.%#,%-G%.%#panic:\ %m,%Ecan't\ load\ package:\ %m,%A%f:%l:%c:\ %m,%A%f:%l:\ %m,%C%*\\s%m,%-G%.%#
setlocal noexpandtab
if &filetype != 'go'
setlocal filetype=go
endif
setlocal fixendofline
setlocal foldcolumn=0
setlocal foldenable
setlocal foldexpr=0
setlocal foldignore=#
setlocal foldlevel=0
setlocal foldmarker={{{,}}}
setlocal foldmethod=manual
setlocal foldminlines=1
setlocal foldnestmax=20
setlocal foldtext=foldtext()
setlocal formatexpr=
setlocal formatoptions=cq
setlocal formatlistpat=^\\s*\\d\\+[\\]:.)}\\t\ ]\\s*
setlocal formatprg=
setlocal grepprg=
setlocal iminsert=2
setlocal imsearch=2
setlocal include=
setlocal includeexpr=
setlocal indentexpr=GoIndent(v:lnum)
setlocal indentkeys=0{,0},:,0#,!^F,o,O,e,<:>,0=},0=)
setlocal noinfercase
setlocal iskeyword=@,48-57,_,192-255
setlocal keywordprg=
setlocal nolinebreak
setlocal nolisp
setlocal lispwords=
setlocal nolist
setlocal makeprg=make
setlocal matchpairs=(:),{:},[:]
setlocal modeline
setlocal modifiable
setlocal nrformats=bin,octal,hex
set number
setlocal number
setlocal numberwidth=4
setlocal omnifunc=gocomplete#Complete
setlocal path=
setlocal nopreserveindent
setlocal nopreviewwindow
setlocal quoteescape=\\
setlocal noreadonly
set relativenumber
setlocal relativenumber
setlocal norightleft
setlocal rightleftcmd=search
setlocal noscrollbind
setlocal shiftwidth=2
setlocal noshortname
setlocal signcolumn=auto
setlocal nosmartindent
setlocal softtabstop=0
setlocal nospell
setlocal spellcapcheck=[.?!]\\_[\\])'\"\	\ ]\\+
setlocal spellfile=
setlocal spelllang=en
setlocal statusline=
setlocal suffixesadd=
setlocal noswapfile
setlocal synmaxcol=3000
if &syntax != 'go'
setlocal syntax=go
endif
setlocal tabstop=2
setlocal tagcase=
setlocal tags=~/go/src/github.com/ilikeorangutans/phts/tags,./tags,./TAGS,tags,TAGS
setlocal textwidth=0
setlocal thesaurus=
setlocal noundofile
setlocal undolevels=-123456
setlocal nowinfixheight
setlocal nowinfixwidth
setlocal wrap
setlocal wrapmargin=0
silent! normal! zE
let s:l = 4 - ((2 * winheight(0) + 9) / 19)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
4
normal! 0
lcd ~/go/src/github.com/ilikeorangutans/phts
wincmd w
argglobal
edit ~/go/src/github.com/ilikeorangutans/phts/test/integration/model/photo_repository_test.go
nnoremap <buffer> <silent>  :call go#def#StackPop(v:count1)
nnoremap <buffer> <silent> ] :call go#def#Jump("split")
nnoremap <buffer> <silent>  :call go#def#Jump("split")
nnoremap <buffer> <silent>  :GoDef
nnoremap <buffer> <silent> K :GoDoc
xnoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('v', 'prev')
onoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('o', 'prev')
nnoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('n', 'prev')
xnoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('v', 'next')
onoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('o', 'next')
nnoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('n', 'next')
xnoremap <buffer> <silent> af :call go#textobj#Function('a')
onoremap <buffer> <silent> af :call go#textobj#Function('a')
let s:cpo_save=&cpo
set cpo&vim
nnoremap <buffer> <silent> g<LeftMouse> <LeftMouse>:GoDef
nnoremap <buffer> <silent> gd :GoDef
xnoremap <buffer> <silent> if :call go#textobj#Function('i')
onoremap <buffer> <silent> if :call go#textobj#Function('i')
nnoremap <buffer> <silent> <C-LeftMouse> <LeftMouse>:GoDef
let &cpo=s:cpo_save
unlet s:cpo_save
setlocal keymap=
setlocal noarabic
setlocal autoindent
setlocal backupcopy=
setlocal balloonexpr=
setlocal nobinary
setlocal nobreakindent
setlocal breakindentopt=
setlocal bufhidden=
setlocal buflisted
setlocal buftype=
setlocal nocindent
setlocal cinkeys=0{,0},0),:,0#,!^F,o,O,e
setlocal cinoptions=
setlocal cinwords=if,else,while,do,for,switch
setlocal colorcolumn=
setlocal comments=s1:/*,mb:*,ex:*/,://
setlocal commentstring=//\ %s
setlocal complete=.,w,b,u,t,i
setlocal concealcursor=
setlocal conceallevel=0
setlocal completefunc=
setlocal nocopyindent
setlocal cryptmethod=
setlocal nocursorbind
setlocal nocursorcolumn
set cursorline
setlocal cursorline
setlocal define=
setlocal dictionary=
setlocal nodiff
setlocal equalprg=
setlocal errorformat=%-G#\ %.%#,%-G%.%#panic:\ %m,%Ecan't\ load\ package:\ %m,%A%f:%l:%c:\ %m,%A%f:%l:\ %m,%C%*\\s%m,%-G%.%#
setlocal noexpandtab
if &filetype != 'go'
setlocal filetype=go
endif
setlocal fixendofline
setlocal foldcolumn=0
setlocal foldenable
setlocal foldexpr=0
setlocal foldignore=#
setlocal foldlevel=0
setlocal foldmarker={{{,}}}
setlocal foldmethod=manual
setlocal foldminlines=1
setlocal foldnestmax=20
setlocal foldtext=foldtext()
setlocal formatexpr=
setlocal formatoptions=cq
setlocal formatlistpat=^\\s*\\d\\+[\\]:.)}\\t\ ]\\s*
setlocal formatprg=
setlocal grepprg=
setlocal iminsert=2
setlocal imsearch=2
setlocal include=
setlocal includeexpr=
setlocal indentexpr=GoIndent(v:lnum)
setlocal indentkeys=0{,0},:,0#,!^F,o,O,e,<:>,0=},0=)
setlocal noinfercase
setlocal iskeyword=@,48-57,_,192-255
setlocal keywordprg=
setlocal nolinebreak
setlocal nolisp
setlocal lispwords=
setlocal nolist
setlocal makeprg=make
setlocal matchpairs=(:),{:},[:]
setlocal modeline
setlocal modifiable
setlocal nrformats=bin,octal,hex
set number
setlocal number
setlocal numberwidth=4
setlocal omnifunc=gocomplete#Complete
setlocal path=
setlocal nopreserveindent
setlocal nopreviewwindow
setlocal quoteescape=\\
setlocal noreadonly
set relativenumber
setlocal relativenumber
setlocal norightleft
setlocal rightleftcmd=search
setlocal noscrollbind
setlocal shiftwidth=2
setlocal noshortname
setlocal signcolumn=auto
setlocal nosmartindent
setlocal softtabstop=0
setlocal nospell
setlocal spellcapcheck=[.?!]\\_[\\])'\"\	\ ]\\+
setlocal spellfile=
setlocal spelllang=en
setlocal statusline=
setlocal suffixesadd=
setlocal noswapfile
setlocal synmaxcol=3000
if &syntax != 'go'
setlocal syntax=go
endif
setlocal tabstop=2
setlocal tagcase=
setlocal tags=~/go/src/github.com/ilikeorangutans/phts/tags,./tags,./TAGS,tags,TAGS
setlocal textwidth=0
setlocal thesaurus=
setlocal noundofile
setlocal undolevels=-123456
setlocal nowinfixheight
setlocal nowinfixwidth
setlocal wrap
setlocal wrapmargin=0
silent! normal! zE
let s:l = 7 - ((2 * winheight(0) + 9) / 19)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
7
normal! 0
lcd ~/go/src/github.com/ilikeorangutans/phts
wincmd w
argglobal
edit ~/go/src/github.com/ilikeorangutans/phts/admin/api/collections.go
nnoremap <buffer> <silent>  :call go#def#StackPop(v:count1)
nnoremap <buffer> <silent> ] :call go#def#Jump("split")
nnoremap <buffer> <silent>  :call go#def#Jump("split")
nnoremap <buffer> <silent>  :GoDef
nnoremap <buffer> <silent> K :GoDoc
xnoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('v', 'prev')
onoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('o', 'prev')
nnoremap <buffer> <silent> [[ :call go#textobj#FunctionJump('n', 'prev')
xnoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('v', 'next')
onoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('o', 'next')
nnoremap <buffer> <silent> ]] :call go#textobj#FunctionJump('n', 'next')
xnoremap <buffer> <silent> af :call go#textobj#Function('a')
onoremap <buffer> <silent> af :call go#textobj#Function('a')
let s:cpo_save=&cpo
set cpo&vim
nnoremap <buffer> <silent> g<LeftMouse> <LeftMouse>:GoDef
nnoremap <buffer> <silent> gd :GoDef
xnoremap <buffer> <silent> if :call go#textobj#Function('i')
onoremap <buffer> <silent> if :call go#textobj#Function('i')
nnoremap <buffer> <silent> <C-LeftMouse> <LeftMouse>:GoDef
let &cpo=s:cpo_save
unlet s:cpo_save
setlocal keymap=
setlocal noarabic
setlocal autoindent
setlocal backupcopy=
setlocal balloonexpr=
setlocal nobinary
setlocal nobreakindent
setlocal breakindentopt=
setlocal bufhidden=
setlocal buflisted
setlocal buftype=
setlocal nocindent
setlocal cinkeys=0{,0},0),:,0#,!^F,o,O,e
setlocal cinoptions=
setlocal cinwords=if,else,while,do,for,switch
setlocal colorcolumn=
setlocal comments=s1:/*,mb:*,ex:*/,://
setlocal commentstring=//\ %s
setlocal complete=.,w,b,u,t,i
setlocal concealcursor=
setlocal conceallevel=0
setlocal completefunc=
setlocal nocopyindent
setlocal cryptmethod=
setlocal nocursorbind
setlocal nocursorcolumn
set cursorline
setlocal cursorline
setlocal define=
setlocal dictionary=
setlocal nodiff
setlocal equalprg=
setlocal errorformat=%-G#\ %.%#,%-G%.%#panic:\ %m,%Ecan't\ load\ package:\ %m,%A%f:%l:%c:\ %m,%A%f:%l:\ %m,%C%*\\s%m,%-G%.%#
setlocal noexpandtab
if &filetype != 'go'
setlocal filetype=go
endif
setlocal fixendofline
setlocal foldcolumn=0
setlocal foldenable
setlocal foldexpr=0
setlocal foldignore=#
setlocal foldlevel=0
setlocal foldmarker={{{,}}}
setlocal foldmethod=manual
setlocal foldminlines=1
setlocal foldnestmax=20
setlocal foldtext=foldtext()
setlocal formatexpr=
setlocal formatoptions=cq
setlocal formatlistpat=^\\s*\\d\\+[\\]:.)}\\t\ ]\\s*
setlocal formatprg=
setlocal grepprg=
setlocal iminsert=2
setlocal imsearch=2
setlocal include=
setlocal includeexpr=
setlocal indentexpr=GoIndent(v:lnum)
setlocal indentkeys=0{,0},:,0#,!^F,o,O,e,<:>,0=},0=)
setlocal noinfercase
setlocal iskeyword=@,48-57,_,192-255
setlocal keywordprg=
setlocal nolinebreak
setlocal nolisp
setlocal lispwords=
setlocal nolist
setlocal makeprg=make
setlocal matchpairs=(:),{:},[:]
setlocal modeline
setlocal modifiable
setlocal nrformats=bin,octal,hex
set number
setlocal number
setlocal numberwidth=4
setlocal omnifunc=gocomplete#Complete
setlocal path=
setlocal nopreserveindent
setlocal nopreviewwindow
setlocal quoteescape=\\
setlocal noreadonly
set relativenumber
setlocal relativenumber
setlocal norightleft
setlocal rightleftcmd=search
setlocal noscrollbind
setlocal shiftwidth=2
setlocal noshortname
setlocal signcolumn=auto
setlocal nosmartindent
setlocal softtabstop=0
setlocal nospell
setlocal spellcapcheck=[.?!]\\_[\\])'\"\	\ ]\\+
setlocal spellfile=
setlocal spelllang=en
setlocal statusline=
setlocal suffixesadd=
setlocal noswapfile
setlocal synmaxcol=3000
if &syntax != 'go'
setlocal syntax=go
endif
setlocal tabstop=2
setlocal tagcase=
setlocal tags=~/go/src/github.com/ilikeorangutans/phts/tags,./tags,./TAGS,tags,TAGS
setlocal textwidth=0
setlocal thesaurus=
setlocal noundofile
setlocal undolevels=-123456
setlocal nowinfixheight
setlocal nowinfixwidth
setlocal wrap
setlocal wrapmargin=0
silent! normal! zE
let s:l = 379 - ((29 * winheight(0) + 30) / 60)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
379
normal! 018|
lcd ~/go/src/github.com/ilikeorangutans/phts
wincmd w
3wincmd w
exe '1resize ' . ((&lines * 20 + 32) / 64)
exe 'vert 1resize ' . ((&columns * 159 + 159) / 319)
exe '2resize ' . ((&lines * 19 + 32) / 64)
exe 'vert 2resize ' . ((&columns * 159 + 159) / 319)
exe '3resize ' . ((&lines * 19 + 32) / 64)
exe 'vert 3resize ' . ((&columns * 159 + 159) / 319)
exe 'vert 4resize ' . ((&columns * 159 + 159) / 319)
tabnext 1
if exists('s:wipebuf')
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20 shortmess=filnxtToO
let s:sx = expand("<sfile>:p:r")."x.vim"
if file_readable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &so = s:so_save | let &siso = s:siso_save
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
