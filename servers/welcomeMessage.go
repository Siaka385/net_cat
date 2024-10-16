package servers

import (
	"strings"
	"time"
)

func WelcomeMessage() string {
	welcomemessage := "Welcome to TCP-Chat!\n" +
		"         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		"|    .       | ' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     -'       --'\n" +
		"[ENTER YOUR NAME]:"

	return welcomemessage
}

func LoadPreviousChats() string {
	prevoiusmessage := strings.Join(UserMessages, "\n")

	return prevoiusmessage
}

func MessageFormat(m string) string {
	times := time.Now()
	formartTime := times.Format("[2006-01-02 15:04:05]")

	return formartTime + "[" + m + "]" + ":"
}
