settings {
    logfacility = "syslog",
    insist = true
}
sync {
      default.rsync,
      source = "/var/lib/eyepi",
      target = "picam@sftp.traitcapture.org:picam/",
      delete = false,
      delay = 260,
      rsync = {
          verbose = true,
          archive = true,
          compress = true,
          dry_run = false,
          _extra = {
            "--remove-source-files",
          }
      }
  }
