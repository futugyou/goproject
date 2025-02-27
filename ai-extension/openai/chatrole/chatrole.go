package chatrole

type ChatRole string

const ChatRoleSystem ChatRole = "system"
const ChatRoleUser ChatRole = "user"
const ChatRoleAssistant ChatRole = "assistant"

var SupportedChatRoles = []ChatRole{ChatRoleSystem, ChatRoleUser, ChatRoleAssistant}
