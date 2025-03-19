# Copyright (C) 2016 Openwrt.org
#
# This is free software, licensed under the Apache License, Version 2.0 .
#
# See /LICENSE for more information.
# This is free software, licensed under the GNU General Public License v2.

include $(TOPDIR)/rules.mk


PKG_NAME:=luci-app-autoreboot
LUCI_TITLE:=LuCI support for Reboot
LUCI_PKGARCH:=all
PKG_VERSION:=1.0.0
PKG_LICENSE:=GPLv3
PKG_LINCESE_FILES:=LICENSE
LUCI_DEPENDS:=+luci-compat

include $(TOPDIR)/feeds/luci/luci.mk

# call BuildPackage - OpenWrt buildroot signature
