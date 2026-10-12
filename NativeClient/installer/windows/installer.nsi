; Erlenwald Online - Windows installer (NSIS, per-user install, no administrator rights needed).
;
; Build (see tools/build_installer.sh):
;   makensis -DVERSION=1.0.6 -DLAUNCHER_DIR=<dir with Launcher.exe + Launcher.pck> [-DGAME_DIR=<exported game, bundles it>]
;            [-DUPDATE_URL=https://updates.example.org] [-DGAME_SERVER=play.example.org:8080] -DOUTFILE=Setup.exe installer.nsi
;
; Without GAME_DIR it is a small "web installer": the launcher downloads the game on its first start.

Unicode true
!include "MUI2.nsh"
!include "FileFunc.nsh"

!ifndef VERSION
  !define VERSION "1.0.6"
!endif
!ifndef OUTFILE
  !define OUTFILE "ErlenwaldOnline-Setup.exe"
!endif
!ifndef UPDATE_URL
  !define UPDATE_URL "https://herrpony.github.io/MMORPG-Prototype/updates"
!endif
!ifndef GAME_SERVER
  !define GAME_SERVER "localhost:8080"
!endif

!define APP "Erlenwald Online"
!define UNKEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\ErlenwaldOnline"

Name "${APP}"
OutFile "${OUTFILE}"
InstallDir "$LOCALAPPDATA\Programs\${APP}"
InstallDirRegKey HKCU "Software\ErlenwaldOnline" "InstallDir"
RequestExecutionLevel user
SetCompressor /SOLID lzma
SetCompressorDictSize 96
BrandingText "${APP} ${VERSION}"
VIProductVersion "${VERSION}.0"
VIAddVersionKey "ProductName" "${APP}"
VIAddVersionKey "FileDescription" "${APP} Installer"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "LegalCopyright" "Erlenwald"

!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN "$INSTDIR\Launcher.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Erlenwald Online starten"
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "German"
!insertmacro MUI_LANGUAGE "English"

Section "Erlenwald Online" SecMain
  SectionIn RO
  SetOutPath "$INSTDIR"
  ; a running launcher cannot be overwritten: ask the user to close it first
  ClearErrors
  IfFileExists "$INSTDIR\Launcher.exe" 0 +2
    Delete "$INSTDIR\Launcher.exe"
  IfErrors 0 +3
    MessageBox MB_OK|MB_ICONEXCLAMATION "Bitte schließe den Erlenwald-Launcher und starte die Installation erneut."
    Abort
  File "${LAUNCHER_DIR}\Launcher.exe"
  File "${LAUNCHER_DIR}\Launcher.pck"
  FileOpen $0 "$INSTDIR\launcher_defaults.json" w
  FileWrite $0 '{"update_url": "${UPDATE_URL}", "server": "${GAME_SERVER}"}'
  FileClose $0
!ifdef GAME_DIR
  SetOutPath "$INSTDIR\game"
  File /r "${GAME_DIR}\*.*"
!endif
  SetOutPath "$INSTDIR"
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  CreateDirectory "$SMPROGRAMS\${APP}"
  CreateShortcut "$SMPROGRAMS\${APP}\${APP}.lnk" "$INSTDIR\Launcher.exe" "" "$INSTDIR\Launcher.exe" 0
  CreateShortcut "$SMPROGRAMS\${APP}\Deinstallieren.lnk" "$INSTDIR\Uninstall.exe"
  CreateShortcut "$DESKTOP\${APP}.lnk" "$INSTDIR\Launcher.exe" "" "$INSTDIR\Launcher.exe" 0
  WriteRegStr HKCU "Software\ErlenwaldOnline" "InstallDir" "$INSTDIR"
  WriteRegStr HKCU "${UNKEY}" "DisplayName" "${APP}"
  WriteRegStr HKCU "${UNKEY}" "DisplayVersion" "${VERSION}"
  WriteRegStr HKCU "${UNKEY}" "Publisher" "Erlenwald"
  WriteRegStr HKCU "${UNKEY}" "InstallLocation" "$INSTDIR"
  WriteRegStr HKCU "${UNKEY}" "DisplayIcon" "$INSTDIR\Launcher.exe"
  WriteRegStr HKCU "${UNKEY}" "UninstallString" '"$INSTDIR\Uninstall.exe"'
  WriteRegDWORD HKCU "${UNKEY}" "NoModify" 1
  WriteRegDWORD HKCU "${UNKEY}" "NoRepair" 1
  ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
  IntFmt $0 "0x%08X" $0
  WriteRegDWORD HKCU "${UNKEY}" "EstimatedSize" "$0"
SectionEnd

Section "Uninstall"
  Delete "$DESKTOP\${APP}.lnk"
  RMDir /r "$SMPROGRAMS\${APP}"
  RMDir /r "$INSTDIR\game"
  RMDir /r "$INSTDIR\.update"
  RMDir /r "$INSTDIR\.partial"
  Delete "$INSTDIR\Launcher.exe"
  Delete "$INSTDIR\Launcher.pck"
  Delete "$INSTDIR\launcher_defaults.json"
  Delete "$INSTDIR\update_launcher.cmd"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  DeleteRegKey HKCU "${UNKEY}"
  DeleteRegKey HKCU "Software\ErlenwaldOnline"
  ; settings and screenshots in %APPDATA%\ErlenwaldOnline stay on purpose
SectionEnd
