
function _() {
    return "{"
}
function __() {
    return "{{"
}
function _n(n) {
    var res = ""
    for(var i = 0; i < n; i++){
        res += "{"
    }
    return res
}

function test(arg) {
    return "js function return :" + arg
}

function test2(arg) {
    return "js function test2"
}
