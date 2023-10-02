$pdf_mode = 1;
@default_files = ('main.tex');
$lualatex = 'lualatex --shell-escape -interaction=nonstopmode -file-line-error -synctex=1 %O %S';
