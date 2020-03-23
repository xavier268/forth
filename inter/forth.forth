
( this file redefines some additionnal forth words )

( =============================================== )
( Compiling forth.forth definitions ) 

: version 0 4 ;     ( -- major minor )     
                    ( marker for preloaded  FORTH definitions )

: .version          ( -- display version major. minor )
    ." Forth version " version swap . ." ." . cr
    ;

( =========== return stack management ================================= )

: ip r@ ;       ( -- ip ) ( provides the ip address, or 0 if high level )

( =========== memory and operations ==================================== )

: +!            ( n addr -- ) ( add n to the cell pointed to by addr )
    dup @       ( n adr x )   
    rot         ( adr x n )   
    + swap !    ( )
    ;     

: ?             ( addr -- ) ( display memory content at provided address )
    @ . ;                   ( memory needs to be allocated/available ! )

: ++ 1 + ;  ( val -- val )  ( increment val )
: -- 1 - ;  ( val -- val )  ( decrement val )

( ============ base management ========================================= )

: decimal 10 base ! ;   ( change base to decimal )
: hex 16 base ! ;       ( change base to hex )
: binary 2 base ! ;     ( change base to binary )

( ============= constants and variables ================================ )

( value -- ) ( variable x : creates x with the initial value )
( -- addr  ) ( x : will put its address on the data stack )
: variable <builds , does>  ;   

( value -- ) ( constant x : creates x with the constant value )
( -- value ) ( x : gets the constant value on stack )
: constant <builds , does> @ ;

( ============== compile =============================================== )
: compile   ( immediately get the cfa of the following word and compile it )
    [ 
        ' ' ,
        ' , ,
    ]
    ; immediate