# envsec Fish shell integration
# Track which vars envsec set, so we can unset them when leaving
set -g __envsec_vars

function __envsec_on_cd --on-variable PWD
    # Unset vars from previous directory
    for var in $__envsec_vars
        set -e $var
    end
    set -g __envsec_vars

    # Load vars for new directory (envsec resolves project + subpath from PWD)
    set -l exports (envsec export --shell fish 2>/dev/null)
    if test $status -eq 0
        for line in $exports
            eval $line
            # Extract var name from "set -gx KEY value" → KEY
            set -a __envsec_vars (string replace -r '^set -gx (\S+) .*' '$1' -- $line)
        end
    end
end

# Run on shell init too (in case we start in a project dir)
__envsec_on_cd
