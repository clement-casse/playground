$pdf_mode = 4;
@default_files = ('main.tex');
$lualatex = 'lualatex -interaction=nonstopmode -file-line-error -synctex=1 -pretex="\pdfvariable suppressoptionalinfo 512\relax" -usepretex --shell-escape %O %S ';
