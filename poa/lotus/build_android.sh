#!/bin/bash

android_sdk_tools_zip="sdk-tools-linux-3859397.zip"
android_ndk_zip="android-ndk-r20-linux-x86_64.zip"

cd ~/binaries/android
curl -so sdk-tools.zip https://dl.google.com/android/repository/$android_sdk_tools_zip
unzip -q sdk-tools.zip
rm sdk-tools.zip
curl -so ndk.zip https://dl.google.com/android/repository/$android_ndk_zip
unzip -q ndk.zip
rm ndk.zip
mv android-ndk-* ndk-bundle
yes|sdkmanager --licenses
sdkmanager "platforms;android-31" "build-tools;32.0.0"
