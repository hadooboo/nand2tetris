#!/bin/bash

cd src
go run *.go ../ArrayTest tokens
go run *.go ../ExpressionLessSquare tokens
go run *.go ../Square tokens
cd ..

../../tools/TextComparer.sh ./ArrayTest/MainT.out.xml ./ArrayTest/MainT.xml

../../tools/TextComparer.sh ./ExpressionLessSquare/MainT.out.xml ./ExpressionLessSquare/MainT.xml
../../tools/TextComparer.sh ./ExpressionLessSquare/SquareT.out.xml ./ExpressionLessSquare/SquareT.xml
../../tools/TextComparer.sh ./ExpressionLessSquare/SquareGameT.out.xml ./ExpressionLessSquare/SquareGameT.xml

../../tools/TextComparer.sh ./Square/MainT.out.xml ./Square/MainT.xml
../../tools/TextComparer.sh ./Square/SquareT.out.xml ./Square/SquareT.xml
../../tools/TextComparer.sh ./Square/SquareGameT.out.xml ./Square/SquareGameT.xml

cd src
go run *.go ../ArrayTest all
go run *.go ../ExpressionLessSquare all
go run *.go ../Square all
cd ..

../../tools/TextComparer.sh ./ArrayTest/Main.out.xml ./ArrayTest/Main.xml

../../tools/TextComparer.sh ./ExpressionLessSquare/Main.out.xml ./ExpressionLessSquare/Main.xml
../../tools/TextComparer.sh ./ExpressionLessSquare/Square.out.xml ./ExpressionLessSquare/Square.xml
../../tools/TextComparer.sh ./ExpressionLessSquare/SquareGame.out.xml ./ExpressionLessSquare/SquareGame.xml

../../tools/TextComparer.sh ./Square/Main.out.xml ./Square/Main.xml
../../tools/TextComparer.sh ./Square/Square.out.xml ./Square/Square.xml
../../tools/TextComparer.sh ./Square/SquareGame.out.xml ./Square/SquareGame.xml
