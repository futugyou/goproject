package main

const (
	EOF                = "\r\n"
	multiLineFileDelim = "_GitHubActionsFileCommandDelimeter_"
	multilineFileCmd   = "%s<<" + multiLineFileDelim + EOF + "%s" + EOF + multiLineFileDelim
)
