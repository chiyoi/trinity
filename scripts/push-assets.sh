cd $(readlink -f $(dirname $0)/..) || return

work() {
    tar -cf tmp assets || return
    az storage blob upload --overwrite -f ./tmp -c assets -n trinity-assets.tar || return
}
work; err=$?
rm -f tmp
return $err
