package internal

import (
	"fmt"
	"os"
)

// PrintUsageWhenConfigFileIsMissing :
func PrintUsageWhenConfigFileIsMissing() {
	fmt.Fprintln(os.Stderr, "if config file is not found. then")
	fmt.Println(`
mkdir -p ~/.config/selfish
cat <<-EOS > ~/.config/selfish/config.json
{
  "access_token": "<your github access token>"
}
EOS
`)
}
