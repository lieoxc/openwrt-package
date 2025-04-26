include $(TOPDIR)/rules.mk

PKG_NAME:=ser2mqtt
PKG_VERSION:=1.0
PKG_RELEASE:=1

# Include mosquitto library as a dependency
PKG_BUILD_DEPENDS:=mosquitto
PKG_LIB_DEPENDS:=libmosquitto

PKG_BUILD_DIR := $(BUILD_DIR)/$(PKG_NAME)

include $(INCLUDE_DIR)/package.mk

define Package/ser2mqtt
	SECTION:=utils
	CATEGORY:=Utilities
	TITLE:=ser2mqtt -- use serial ttyS1 to mqtt
	DEPENDS:=+libmosquitto
endef

define Package/ser2mqtt/description
    If you can't figure out what this program does, you're probably
    brain-dead and need immediate medical attention.
endef

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	$(CP) ./src/* $(PKG_BUILD_DIR)/
endef

define Build/Configure
endef

define Build/Compile
	$(MAKE) -C $(PKG_BUILD_DIR) \
		CC="$(TARGET_CC)" \
		CFLAGS="$(TARGET_CFLAGS) -Wall" \
		LDFLAGS="$(TARGET_LDFLAGS) -lmosquitto"
endef

define Package/ser2mqtt/install
	$(INSTALL_DIR) $(1)/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/ser2mqtt $(1)/bin/
	$(CP) ./files/* $(1)/
endef

$(eval $(call BuildPackage,ser2mqtt))
