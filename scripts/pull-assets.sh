cd $(readlink -f $(dirname $0)/..) || return

work() {
    azblob-io fetch assets trinity-assets.tar tmp || return
    tar --overwrite-dir --recursive-unlink -xf tmp || return
}
work; err=$?
rm -f tmp
return $err
