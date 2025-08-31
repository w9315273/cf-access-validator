#
# Copyright (C) 2025 w9315273
#
# This is free software, licensed under the Apache License, Version 2.0.
#

include $(TOPDIR)/rules.mk

PKG_NAME:=cf-auth
PKG_VERSION:=0.1.0

PKG_LICENSE:=MIT
PKG_LICENSE_FILES:=LICENSE
PKG_MAINTAINER:=w9315273

PKG_BUILD_DEPENDS:=golang/host

GO_MODULE_PATH:=github.com/w9315273/cf-access-validator
GO_PKG:=$(GO_MODULE_PATH)
GO_PKG_BUILD_PKG:=$(GO_MODULE_PATH)/apps/cf-auth
GO_PKG_LDFLAGS:=-s -w

include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk

define Package/cf-auth
  SECTION:=net
  CATEGORY:=Network
  TITLE:=Cloudflare Access JWT validator
  URL:=https://github.com/w9315273/cf-access-validator
  DEPENDS:=$(GO_ARCH_DEPENDS)
endef

define Package/cf-auth/description
Minimal Cloudflare Access JWT validator for OpenWrt.
endef

define Build/Prepare
	$(call GoPackage/Build/Prepare)
	$(CP) $(CURDIR)/src/* $(PKG_BUILD_DIR)/
endef

define Package/cf-auth/install
	$(call GoPackage/Package/Install/Bin,$(1))
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) $(CURDIR)/src/packaging/openwrt/files/etc/init.d/cf-auth $(1)/etc/init.d/cf-auth
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) $(CURDIR)/src/packaging/openwrt/files/etc/config/cf-auth $(1)/etc/config/cf-auth
endef

$(eval $(call GoBinPackage,$(PKG_NAME)))
$(eval $(call BuildPackage,$(PKG_NAME)))