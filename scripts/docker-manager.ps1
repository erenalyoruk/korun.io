#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Docker Compose Manager for Korun.io Development Environment
.DESCRIPTION
    Manages Docker Compose services with various operations like start, stop, rebuild, reset, etc.
.PARAMETER Action
    The action to perform (up, down, restart, rebuild, reset, logs, status, migrate, clean)
.PARAMETER Service
    Optional service name to target specific service
.PARAMETER Follow
    Follow logs in real-time (only for logs action)
.PARAMETER SkipMigration
    Skip automatic migration after rebuild/reset operations
.EXAMPLE
    .\docker-manager.ps1 up
    .\docker-manager.ps1 rebuild secret-service-dev
    .\docker-manager.ps1 logs -Follow
    .\docker-manager.ps1 reset
    .\docker-manager.ps1 rebuild -SkipMigration
#>

param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("up", "down", "restart", "rebuild", "reset", "logs", "status", "migrate", "clean", "fix-networks", "help")]
    [string]$Action,

    [Parameter(Mandatory = $false)]
    [string]$Service = "",

    [Parameter(Mandatory = $false)]
    [switch]$Follow,

    [Parameter(Mandatory = $false)]
    [switch]$SkipMigration
)

# Configuration
$ComposeFile = "docker-compose.dev.yml"
$ComposePath = "infrastructure/docker/development"
$EnvFile = ".env"

# Colors for output
$Red = [System.ConsoleColor]::Red
$Green = [System.ConsoleColor]::Green
$Yellow = [System.ConsoleColor]::Yellow
$Blue = [System.ConsoleColor]::Blue
$Cyan = [System.ConsoleColor]::Cyan

function Write-ColoredOutput {
    param(
        [string]$Message,
        [System.ConsoleColor]$Color = [System.ConsoleColor]::White
    )
    Write-Host $Message -ForegroundColor $Color
}

function Test-DockerCompose {
    try {
        docker-compose --version | Out-Null
        return $true
    }
    catch {
        Write-ColoredOutput "❌ Docker Compose not found. Please install Docker Desktop." $Red
        return $false
    }
}

function Test-ComposeFile {
    $fullPath = Join-Path $ComposePath $ComposeFile
    if (-not (Test-Path $fullPath)) {
        Write-ColoredOutput "❌ Compose file not found: $fullPath" $Red
        return $false
    }
    return $true
}

function Get-ComposeCommand {
    param([string]$cmd)
    $fullPath = Join-Path $ComposePath $ComposeFile
    return "docker-compose -f `"$fullPath`" $cmd"
}

function Invoke-ComposeCommand {
    param([string]$cmd)
    $fullCommand = Get-ComposeCommand $cmd
    Write-ColoredOutput "🔄 Executing: $fullCommand" $Blue
    Invoke-Expression $fullCommand
}

function Show-Help {
    Write-ColoredOutput "🐳 Docker Compose Manager for Korun.io" $Cyan
    Write-ColoredOutput ("=" * 50) $Cyan
    Write-ColoredOutput ""
    Write-ColoredOutput "Available Actions:" $Green
    Write-ColoredOutput "  up          - Start all services" $Yellow
    Write-ColoredOutput "  down        - Stop all services" $Yellow
    Write-ColoredOutput "  restart     - Restart services" $Yellow
    Write-ColoredOutput "  rebuild     - Rebuild and restart services (auto-runs migration)" $Yellow
    Write-ColoredOutput "  reset       - Nuclear reset (remove everything and rebuild, auto-runs migration)" $Yellow
    Write-ColoredOutput "  logs        - Show logs (use -Follow for real-time)" $Yellow
    Write-ColoredOutput "  status      - Show service status" $Yellow
    Write-ColoredOutput "  migrate     - Run database migrations" $Yellow
    Write-ColoredOutput "  clean       - Clean up unused Docker resources" $Yellow
    Write-ColoredOutput "  help        - Show this help" $Yellow
    Write-ColoredOutput ""
    Write-ColoredOutput "Parameters:" $Green
    Write-ColoredOutput "  -Service    - Target specific service" $Yellow
    Write-ColoredOutput "  -Follow     - Follow logs in real-time" $Yellow
    Write-ColoredOutput "  -SkipMigration - Skip automatic migration after rebuild/reset" $Yellow
    Write-ColoredOutput ""
    Write-ColoredOutput "Examples:" $Green
    Write-ColoredOutput "  .\docker-manager.ps1 up" $Cyan
    Write-ColoredOutput "  .\docker-manager.ps1 rebuild secret-service-dev" $Cyan
    Write-ColoredOutput "  .\docker-manager.ps1 logs -Follow" $Cyan
    Write-ColoredOutput "  .\docker-manager.ps1 reset" $Cyan
    Write-ColoredOutput "  .\docker-manager.ps1 rebuild -SkipMigration" $Cyan
}

function Start-Services {
    Write-ColoredOutput "🚀 Starting services..." $Green
    if ($Service) {
        Invoke-ComposeCommand "up -d $Service"
    } else {
        Invoke-ComposeCommand "up -d"
    }
    Write-ColoredOutput "✅ Services started!" $Green
}

function Stop-Services {
    Write-ColoredOutput "🛑 Stopping services..." $Yellow
    if ($Service) {
        Invoke-ComposeCommand "stop $Service"
    } else {
        Invoke-ComposeCommand "down"
    }
    Write-ColoredOutput "✅ Services stopped!" $Green
}

function Restart-Services {
    Write-ColoredOutput "🔄 Restarting services..." $Yellow
    if ($Service) {
        Invoke-ComposeCommand "restart $Service"
    } else {
        Invoke-ComposeCommand "restart"
    }
    Write-ColoredOutput "✅ Services restarted!" $Green
}

function Rebuild-Services {
    Write-ColoredOutput "🔨 Rebuilding services..." $Yellow
    if ($Service) {
        Invoke-ComposeCommand "build --no-cache $Service"
        Invoke-ComposeCommand "up -d $Service"
    } else {
        Invoke-ComposeCommand "build --no-cache"
        Invoke-ComposeCommand "up -d"
    }
    Write-ColoredOutput "✅ Services rebuilt and started!" $Green

    # Auto-run migration unless explicitly skipped
    if (-not $SkipMigration) {
        Write-ColoredOutput "🔄 Auto-running migrations after rebuild..." $Blue
        Run-Migration
    } else {
        Write-ColoredOutput "⚠️ Skipping migration as requested" $Yellow
    }
}

function Reset-Everything {
    Write-ColoredOutput "💥 Nuclear reset - This will remove everything!" $Red
    $confirm = Read-Host "Are you sure? This will delete all data! (y/N)"

    if ($confirm -eq 'y' -or $confirm -eq 'Y') {
        Write-ColoredOutput "🗑️ Stopping and removing everything..." $Yellow
        Invoke-ComposeCommand "down --rmi all --volumes --remove-orphans"

        Write-ColoredOutput "🔨 Rebuilding from scratch..." $Yellow
        Invoke-ComposeCommand "build --no-cache"

        Write-ColoredOutput "🚀 Starting fresh services..." $Green
        Invoke-ComposeCommand "up -d"

        Write-ColoredOutput "✅ Complete reset finished!" $Green

        # Auto-run migration unless explicitly skipped
        if (-not $SkipMigration) {
            Write-ColoredOutput "🔄 Auto-running migrations after reset..." $Blue
            Run-Migration
        } else {
            Write-ColoredOutput "⚠️ Skipping migration as requested" $Yellow
        }
    } else {
        Write-ColoredOutput "❌ Reset cancelled." $Yellow
    }
}

function Show-Logs {
    Write-ColoredOutput "📋 Showing logs..." $Blue
    if ($Follow) {
        if ($Service) {
            Invoke-ComposeCommand "logs -f $Service"
        } else {
            Invoke-ComposeCommand "logs -f"
        }
    } else {
        if ($Service) {
            Invoke-ComposeCommand "logs --tail=50 $Service"
        } else {
            Invoke-ComposeCommand "logs --tail=50"
        }
    }
}

function Show-Status {
    Write-ColoredOutput "📊 Service Status:" $Blue
    Invoke-ComposeCommand "ps"
    Write-ColoredOutput ""
    Write-ColoredOutput "🖼️ Images:" $Blue
    Invoke-ComposeCommand "images"
}

function Run-Migration {
    Write-ColoredOutput "🗄️ Running database migrations..." $Blue

    # Wait for services to be fully ready
    Write-ColoredOutput "⏳ Waiting for services to be ready..." $Blue
    Start-Sleep -Seconds 3

    # Check network status and ensure it's properly established
    Write-ColoredOutput "🔗 Verifying network connectivity..." $Blue
    $networkName = "development_korun-io-network"  # Adjust this to match your network name

    try {
        # Check if network exists and is active
        $networkInfo = docker network inspect $networkName 2>$null
        if (-not $networkInfo) {
            Write-ColoredOutput "⚠️ Network not found, recreating services..." $Yellow
            Invoke-ComposeCommand "down"
            Invoke-ComposeCommand "up -d"
            Start-Sleep -Seconds 5
        }

        # Ensure migration container can access the network
        Write-ColoredOutput "🚀 Starting migration with network verification..." $Blue
        Invoke-ComposeCommand "run --rm migrate-dev"
        Write-ColoredOutput "✅ Migration completed!" $Green
    }
    catch {
        Write-ColoredOutput "⚠️ Migration failed. Trying network reset..." $Yellow

        try {
            # Stop everything and restart to reset network state
            Write-ColoredOutput "🔄 Resetting network state..." $Blue
            Invoke-ComposeCommand "down"
            Start-Sleep -Seconds 2
            Invoke-ComposeCommand "up -d"
            Start-Sleep -Seconds 5
            Invoke-ComposeCommand "run --rm migrate-dev"
            Write-ColoredOutput "✅ Migration completed after network reset!" $Green
        }
        catch {
            Write-ColoredOutput "❌ Migration failed. Please run manually: .\docker-manager.ps1 migrate" $Red
            Write-ColoredOutput "Or try: .\docker-manager.ps1 fix-networks" $Yellow
        }
    }
}

function Clean-Docker {
    Write-ColoredOutput "🧹 Cleaning up Docker resources..." $Yellow

    Write-ColoredOutput "Stopping all containers..." $Blue
    Invoke-ComposeCommand "down"

    Write-ColoredOutput "Removing unused containers..." $Blue
    docker container prune -f

    Write-ColoredOutput "Removing unused images..." $Blue
    docker image prune -a -f

    Write-ColoredOutput "Removing unused volumes..." $Blue
    docker volume prune -f

    Write-ColoredOutput "Removing unused networks..." $Blue
    docker network prune -f

    Write-ColoredOutput "✅ Docker cleanup completed!" $Green
}

function Fix-Networks {
    Write-ColoredOutput "🔧 Fixing network issues..." $Yellow

    Write-ColoredOutput "Stopping services..." $Blue
    Invoke-ComposeCommand "down"

    Write-ColoredOutput "Cleaning networks..." $Blue
    docker network prune -f

    Write-ColoredOutput "Recreating services..." $Blue
    Invoke-ComposeCommand "up -d --force-recreate"

    Write-ColoredOutput "✅ Network issues fixed!" $Green
}

# Main execution
if ($Action -eq "help") {
    Show-Help
    exit 0
}

# Validate requirements
if (-not (Test-DockerCompose)) {
    exit 1
}

if (-not (Test-ComposeFile)) {
    exit 1
}

# Check if we're in the right directory
if (-not (Test-Path $ComposePath)) {
    Write-ColoredOutput "❌ Please run this script from the project root directory." $Red
    exit 1
}

# Load environment variables if .env exists
$envPath = Join-Path $ComposePath $EnvFile
if (Test-Path $envPath) {
    Write-ColoredOutput "📄 Loading environment variables from $envPath" $Blue
    Get-Content $envPath | ForEach-Object {
        if ($_ -match '^([^=]+)=(.*)$') {
            [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2])
        }
    }
}

# Execute the requested action
switch ($Action) {
    "up" { Start-Services }
    "down" { Stop-Services }
    "restart" { Restart-Services }
    "rebuild" { Rebuild-Services }
    "reset" { Reset-Everything }
    "logs" { Show-Logs }
    "status" { Show-Status }
    "migrate" { Run-Migration }
    "clean" { Clean-Docker }
    "fix-networks" { Fix-Networks }
    default {
        Write-ColoredOutput "❌ Unknown action: $Action" $Red
        Show-Help
        exit 1
    }
}
