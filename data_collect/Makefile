include $(TOPDIR)/rules.mk

PKG_NAME:=data_collect
PKG_RELEASE:=1

PKG_BUILD_DIR := $(BUILD_DIR)/$(PKG_NAME)

include $(INCLUDE_DIR)/package.mk

define Package/data_collect
	SECTION:=utils
	CATEGORY:=Utilities
	TITLE:=data_collect  -- use ttyS1 send Modbus RTU data,and recv Modbus RTU data. Pub to Mqtt borker
endef

define Package/data_collect/description
    If you can't figure out what this program does, you're probably
    brain-dead and need immediate medical attention.
endef

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
endef

define Build/Configure
endef

define Build/Compile
endef

define Package/data_collect/install
	$(INSTALL_DIR) $(1)/bin
	$(INSTALL_BIN) ./bin/data_collect $(1)/bin/
	$(CP) ./files/* $(1)/
endef

$(eval $(call BuildPackage,data_collect))
