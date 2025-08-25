package utils

// Standard error codes used in JSON error responses.
const (
	CodeEmptyBody       = "EMPTY_BODY"
	CodeJSONSyntax      = "JSON_SYNTAX"
	CodeJSONType        = "JSON_TYPE"
	CodeUnknownField    = "UNKNOWN_FIELD"
	CodeInvalidJSON     = "INVALID_JSON"
	CodeMultipleObjects = "MULTIPLE_OBJECTS"
	CodeMissingFields   = "MISSING_FIELDS"
	CodeInvalidLanguage = "INVALID_LANGUAGE"
	CodeDuplicatePack   = "DUPLICATE_PACK"
	CodeInvalidPack     = "INVALID_PACK"
	CodeDuplicateVocab  = "DUPLICATE_VOCAB"
	CodeInvalidVocab    = "INVALID_VOCAB"
	CodeInvalidPacks    = "INVALID_PACKS"
	CodeInvalidFileType = "INVALID_FILE_TYPE"
	CodeFileTooLarge    = "FILE_TOO_LARGE"
	CodeInternal        = "INTERNAL"
)
