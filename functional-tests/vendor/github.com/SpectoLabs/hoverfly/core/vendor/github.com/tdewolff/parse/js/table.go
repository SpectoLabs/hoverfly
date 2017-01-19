package js // import "github.com/tdewolff/parse/js"

var regexpStateByte = [128]bool{
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false, //
	true,  // !
	false, // "
	false, // #
	false, // $
	false, // %
	true,  // &
	false, // '
	true,  // (
	false, // )
	true,  // *
	true,  // +
	true,  // ,
	true,  // -
	false, // .
	true,  // /
	false, // 0
	false, // 1
	false, // 2
	false, // 3
	false, // 4
	false, // 5
	false, // 6
	false, // 7
	false, // 8
	false, // 9
	true,  // :
	true,  // ;
	false, // <
	true,  // =
	false, // >
	true,  // ?
	false, // @
	false, // A-Z...
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	true,  // [
	false, // \
	false, // ]
	false, // ^
	false, // _
	false, // `
	false, // a-z...
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	true, // {
	true, // |
	true, // }
	true, // ~
	false,
}
