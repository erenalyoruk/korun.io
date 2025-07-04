@echo off
REM Docker Compose Manager Batch Wrapper
REM This calls the PowerShell script with the same arguments

pwsh -ExecutionPolicy Bypass -File "%~dp0./scripts/docker-manager.ps1" %*
