# This file was generated by GoReleaser. DO NOT EDIT.
class GitSemver < Formula
  desc ""
  homepage ""
  url "https://github.com/meinto/git-semver/releases/download/v5.0.0/semver_5.0.0_darwin_x86_64.tar.gz"
  version "5.0.0"
  sha256 "5712584ab226a23851e92dca9f7a367be6a604a49d4856b5ce752639290b6b7f"
  
  depends_on "git"

  def install
    bin.install "git-semver"
  end
end
