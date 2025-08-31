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
GO_PKG:=$(GO_MODULE_PATH)/apps/cf-auth
GO_PKG_BUILD_PKG:=$(GO_PKG)
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
	$(INSTALL_DIR) $(PKG_BUILD_DIR)/src/$(GO_MODULE_PATH)
	$(CP) -a $(CURDIR)/apps $(PKG_BUILD_DIR)/src/$(GO_MODULE_PATH)/
	$(CP) -a $(CURDIR)/go.mod $(CURDIR)/go.sum $(PKG_BUILD_DIR)/src/$(GO_MODULE_PATH)/
endef

define Build/Compile
	( \
		cd $(PKG_BUILD_DIR)/src/$(GO_MODULE_PATH)/apps/cf-auth ; \
		GO111MODULE=on CGO_ENABLED=0 \
		$(STAGING_DIR_HOSTPKG)/bin/go build -trimpath -ldflags "$(GO_PKG_LDFLAGS)" \
			-o $(PKG_BUILD_DIR)/bin/cf-auth ; \
	)
endef

define Package/cf-auth/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/bin/cf-auth $(1)/usr/bin/cf-auth

	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) $(CURDIR)/files/etc/init.d/cf-auth $(1)/etc/init.d/cf-auth

	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) $(CURDIR)/files/etc/config/cf-auth $(1)/etc/config/cf-auth
endef

$(eval $(call BuildPackage,$(PKG_NAME)))