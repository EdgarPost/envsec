package shell

import "fmt"

func FishHook(envsecPath string) string {
	return fmt.Sprintf(`# envsec Fish shell integration
set -g __envsec_vars

function __envsec_load
    # Unset vars from previous load
    for var in $__envsec_vars
        set -e $var
    end
    set -g __envsec_vars

    # Load vars for current directory
    set -l exports (%[1]s export --shell fish 2>/dev/null)
    if test $status -eq 0
        for line in $exports
            eval $line
            set -a __envsec_vars (string replace -r '^set -gx (\S+) .*' '$1' -- $line)
        end
    end
end

function __envsec_on_cd --on-variable PWD
    __envsec_load
end

# Wrapper: reload env after mutations
function envsec --wraps=%[1]s
    %[1]s $argv
    set -l cmd_status $status
    switch "$argv[1]"
        case set rm import init
            __envsec_load
    end
    return $cmd_status
end

# Load on shell init
__envsec_load
`, envsecPath)
}
