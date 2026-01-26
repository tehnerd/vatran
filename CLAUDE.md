this projects contains c bindings for katran load balancer
main goal of this bindings is to be able to link katranLb library against go via cgo
the katrans library source code could be found in upper dir from this project
(e.g. this git project dir is <pdir>; katran is <pdir>/../lib)
include/ contains public interfaces
src/ counts implimintation source code

all code in this project must be properly documented (public methods) w/ docstring format for 
all input param

Do not use/read any GO code outside of this project directory

Do not add Facebook's copyrights at the top of the files

main language of the project is
1. C (for katran's library bindings)
2. Go (for the main backend code)
