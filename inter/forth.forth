
( this file redefines some additionnal forth words )

( =============================================== )
( Compiling forth.forth definitions ) 

: version 0 4 ;     ( -- major minor )     
                    ( marker for preloaded  FORTH definitions )

: .version          ( -- display version major. minor )
    ." Forth version " version swap . ." ." . cr
    ;


: +!            ( n addr -- ) ( add n to the cell pointed to by addr )
    dup @       ( n adr x )   
    rot         ( adr x n )   
    + swap !    ( )
    ;     

( value -- ) ( variable x : creates x with the initial value )
( -- addr  ) ( x : will put its address on the data stack )
: variable <builds , does>  ;   

( value -- ) ( constant x : creates x with the constant value )
( -- value ) ( x : gets the constant value on stack )
: constant <builds , does> @ ;

