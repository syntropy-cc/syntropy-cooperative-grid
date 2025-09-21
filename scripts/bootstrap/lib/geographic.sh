#!/bin/bash

# Syntropy Cooperative Grid - Geographic Detection System
# Version: 2.0.0

# Enhanced geographic coordinate detection with multiple fallback methods
detect_coordinates() {
    local manual_coords="$1"
    
    if [ -n "$manual_coords" ]; then
        echo "$manual_coords:manual:Manual:Entry"
        return 0
    fi
    
    log INFO "Detecting geographic coordinates using multiple methods..."
    
    # Method 1: ipapi.co (most accurate)
    if attempt_ipapi_detection; then
        return 0
    fi
    
    # Method 2: ipinfo.io (backup)
    if attempt_ipinfo_detection; then
        return 0
    fi
    
    # Method 3: ip-api.com (another backup)
    if attempt_ipapi_com_detection; then
        return 0
    fi
    
    # Method 4: Enhanced timezone-based approximation
    timezone_based_detection
}

# Attempt detection using ipapi.co
attempt_ipapi_detection() {
    log DEBUG "Attempting coordinate detection via ipapi.co..."
    
    local result=$(timeout 15 curl -s "http://ipapi.co/json" 2>/dev/null || true)
    if [ -n "$result" ]; then
        local coords_result=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    lat = data.get('latitude')
    lon = data.get('longitude')
    city = data.get('city', 'unknown')
    country = data.get('country_name', 'unknown')
    if lat and lon and str(lat) != 'None' and str(lon) != 'None':
        print(f'{lat},{lon}:ip_geolocation_ipapi:{city}:{country}')
        exit(0)
except Exception:
    pass
exit(1)
" 2>/dev/null)
        
        if [ $? -eq 0 ] && [ -n "$coords_result" ] && [[ "$coords_result" != *"None"* ]]; then
            echo "$coords_result"
            log SUCCESS "Coordinates detected via ipapi.co"
            return 0
        fi
    fi
    
    log DEBUG "ipapi.co detection failed"
    return 1
}

# Attempt detection using ipinfo.io
attempt_ipinfo_detection() {
    log DEBUG "Attempting coordinate detection via ipinfo.io..."
    
    local result=$(timeout 15 curl -s "http://ipinfo.io/json" 2>/dev/null || true)
    if [ -n "$result" ]; then
        local coords_result=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    loc = data.get('loc', '')
    city = data.get('city', 'unknown')
    country = data.get('country', 'unknown')
    if loc and ',' in loc:
        print(f'{loc}:ip_geolocation_ipinfo:{city}:{country}')
        exit(0)
except Exception:
    pass
exit(1)
" 2>/dev/null)
        
        if [ $? -eq 0 ] && [ -n "$coords_result" ] && [[ "$coords_result" == *","* ]]; then
            echo "$coords_result"
            log SUCCESS "Coordinates detected via ipinfo.io"
            return 0
        fi
    fi
    
    log DEBUG "ipinfo.io detection failed"
    return 1
}

# Attempt detection using ip-api.com
attempt_ipapi_com_detection() {
    log DEBUG "Attempting coordinate detection via ip-api.com..."
    
    local result=$(timeout 15 curl -s "http://ip-api.com/json" 2>/dev/null || true)
    if [ -n "$result" ]; then
        local coords_result=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    lat = data.get('lat')
    lon = data.get('lon')
    city = data.get('city', 'unknown')
    country = data.get('country', 'unknown')
    if lat and lon and str(lat) != 'None' and str(lon) != 'None':
        print(f'{lat},{lon}:ip_geolocation_ipapi_com:{city}:{country}')
        exit(0)
except Exception:
    pass
exit(1)
" 2>/dev/null)
        
        if [ $? -eq 0 ] && [ -n "$coords_result" ] && [[ "$coords_result" != *"None"* ]]; then
            echo "$coords_result"
            log SUCCESS "Coordinates detected via ip-api.com"
            return 0
        fi
    fi
    
    log DEBUG "ip-api.com detection failed"
    return 1
}

# Enhanced timezone-based approximation with major cities
timezone_based_detection() {
    log DEBUG "Using timezone-based coordinate approximation..."
    
    local timezone=$(timedatectl show --property=Timezone --value 2>/dev/null || echo "UTC")
    log DEBUG "Detected timezone: $timezone"
    
    case "$timezone" in
        # Brazil - More comprehensive coverage
        "America/Sao_Paulo"|"America/Bahia"|"America/Fortaleza"|"America/Recife") 
            echo "-23.5505,-46.6333:timezone_brazil:São Paulo:Brazil" ;;
        "America/Manaus") 
            echo "-3.1190,-60.0217:timezone_brazil:Manaus:Brazil" ;;
        "America/Rio_Branco") 
            echo "-9.9750,-67.8243:timezone_brazil:Rio Branco:Brazil" ;;
        "America/Belem") 
            echo "-1.4558,-48.5044:timezone_brazil:Belém:Brazil" ;;
        
        # North America - Extended coverage
        "America/New_York"|"America/Detroit"|"America/Toronto") 
            echo "40.7128,-74.0060:timezone_us_eastern:New York:United States" ;;
        "America/Chicago"|"America/Winnipeg") 
            echo "41.8781,-87.6298:timezone_us_central:Chicago:United States" ;;
        "America/Denver"|"America/Edmonton") 
            echo "39.7392,-104.9903:timezone_us_mountain:Denver:United States" ;;
        "America/Los_Angeles"|"America/Vancouver") 
            echo "34.0522,-118.2437:timezone_us_pacific:Los Angeles:United States" ;;
        "America/Mexico_City") 
            echo "19.4326,-99.1332:timezone_mexico:Mexico City:Mexico" ;;
        "America/Anchorage") 
            echo "61.2181,-149.9003:timezone_us_alaska:Anchorage:United States" ;;
        "America/Phoenix") 
            echo "33.4484,-112.0740:timezone_us_arizona:Phoenix:United States" ;;
        
        # Europe - Comprehensive coverage
        "Europe/London"|"Europe/Belfast") 
            echo "51.5074,-0.1278:timezone_uk:London:United Kingdom" ;;
        "Europe/Paris") 
            echo "48.8566,2.3522:timezone_france:Paris:France" ;;
        "Europe/Berlin") 
            echo "52.5200,13.4050:timezone_germany:Berlin:Germany" ;;
        "Europe/Madrid") 
            echo "40.4168,-3.7038:timezone_spain:Madrid:Spain" ;;
        "Europe/Rome") 
            echo "41.9028,12.4964:timezone_italy:Rome:Italy" ;;
        "Europe/Amsterdam") 
            echo "52.3676,4.9041:timezone_netherlands:Amsterdam:Netherlands" ;;
        "Europe/Stockholm") 
            echo "59.3293,18.0686:timezone_sweden:Stockholm:Sweden" ;;
        "Europe/Oslo") 
            echo "59.9139,10.7522:timezone_norway:Oslo:Norway" ;;
        "Europe/Copenhagen") 
            echo "55.6761,12.5683:timezone_denmark:Copenhagen:Denmark" ;;
        "Europe/Helsinki") 
            echo "60.1699,24.9384:timezone_finland:Helsinki:Finland" ;;
        "Europe/Warsaw") 
            echo "52.2297,21.0122:timezone_poland:Warsaw:Poland" ;;
        "Europe/Prague") 
            echo "50.0755,14.4378:timezone_czech:Prague:Czech Republic" ;;
        "Europe/Vienna") 
            echo "48.2082,16.3738:timezone_austria:Vienna:Austria" ;;
        "Europe/Zurich") 
            echo "47.3769,8.5417:timezone_switzerland:Zurich:Switzerland" ;;
        "Europe/Brussels") 
            echo "50.8503,4.3517:timezone_belgium:Brussels:Belgium" ;;
        "Europe/Dublin") 
            echo "53.3498,-6.2603:timezone_ireland:Dublin:Ireland" ;;
        "Europe/Lisbon") 
            echo "38.7223,-9.1393:timezone_portugal:Lisbon:Portugal" ;;
        "Europe/Athens") 
            echo "37.9838,23.7275:timezone_greece:Athens:Greece" ;;
        "Europe/Budapest") 
            echo "47.4979,19.0402:timezone_hungary:Budapest:Hungary" ;;
        "Europe/Bucharest") 
            echo "44.4268,26.1025:timezone_romania:Bucharest:Romania" ;;
        "Europe/Moscow") 
            echo "55.7558,37.6176:timezone_russia:Moscow:Russia" ;;
        
        # Asia - Extensive coverage
        "Asia/Tokyo") 
            echo "35.6762,139.6503:timezone_japan:Tokyo:Japan" ;;
        "Asia/Seoul") 
            echo "37.5665,126.9780:timezone_korea:Seoul:South Korea" ;;
        "Asia/Shanghai"|"Asia/Beijing") 
            echo "31.2304,121.4737:timezone_china:Shanghai:China" ;;
        "Asia/Hong_Kong") 
            echo "22.3193,114.1694:timezone_hongkong:Hong Kong:Hong Kong" ;;
        "Asia/Singapore") 
            echo "1.3521,103.8198:timezone_singapore:Singapore:Singapore" ;;
        "Asia/Mumbai"|"Asia/Kolkata"|"Asia/Calcutta") 
            echo "19.0760,72.8777:timezone_india:Mumbai:India" ;;
        "Asia/Dubai") 
            echo "25.2048,55.2708:timezone_uae:Dubai:United Arab Emirates" ;;
        "Asia/Tehran") 
            echo "35.6892,51.3890:timezone_iran:Tehran:Iran" ;;
        "Asia/Bangkok") 
            echo "13.7563,100.5018:timezone_thailand:Bangkok:Thailand" ;;
        "Asia/Jakarta") 
            echo "-6.2088,106.8456:timezone_indonesia:Jakarta:Indonesia" ;;
        "Asia/Manila") 
            echo "14.5995,120.9842:timezone_philippines:Manila:Philippines" ;;
        "Asia/Kuala_Lumpur") 
            echo "3.1390,101.6869:timezone_malaysia:Kuala Lumpur:Malaysia" ;;
        "Asia/Ho_Chi_Minh") 
            echo "10.8231,106.6297:timezone_vietnam:Ho Chi Minh City:Vietnam" ;;
        "Asia/Taipei") 
            echo "25.0330,121.5654:timezone_taiwan:Taipei:Taiwan" ;;
        
        # Australia & Oceania
        "Australia/Sydney") 
            echo "-33.8688,151.2093:timezone_australia:Sydney:Australia" ;;
        "Australia/Melbourne") 
            echo "-37.8136,144.9631:timezone_australia:Melbourne:Australia" ;;
        "Australia/Brisbane") 
            echo "-27.4698,153.0251:timezone_australia:Brisbane:Australia" ;;
        "Australia/Perth") 
            echo "-31.9505,115.8605:timezone_australia:Perth:Australia" ;;
        "Australia/Adelaide") 
            echo "-34.9285,138.6007:timezone_australia:Adelaide:Australia" ;;
        "Pacific/Auckland") 
            echo "-36.8485,174.7633:timezone_newzealand:Auckland:New Zealand" ;;
        "Pacific/Fiji") 
            echo "-18.1248,178.4501:timezone_fiji:Suva:Fiji" ;;
        
        # Africa
        "Africa/Johannesburg") 
            echo "-26.2041,28.0473:timezone_southafrica:Johannesburg:South Africa" ;;
        "Africa/Cairo") 
            echo "30.0444,31.2357:timezone_egypt:Cairo:Egypt" ;;
        "Africa/Lagos") 
            echo "6.5244,3.3792:timezone_nigeria:Lagos:Nigeria" ;;
        "Africa/Nairobi") 
            echo "-1.2921,36.8219:timezone_kenya:Nairobi:Kenya" ;;
        "Africa/Casablanca") 
            echo "33.5731,-7.5898:timezone_morocco:Casablanca:Morocco" ;;
        "Africa/Algiers") 
            echo "36.7538,3.0588:timezone_algeria:Algiers:Algeria" ;;
        
        # Middle East
        "Asia/Jerusalem") 
            echo "31.7683,35.2137:timezone_israel:Jerusalem:Israel" ;;
        "Asia/Riyadh") 
            echo "24.7136,46.6753:timezone_saudi:Riyadh:Saudi Arabia" ;;
        "Asia/Kuwait") 
            echo "29.3759,47.9774:timezone_kuwait:Kuwait City:Kuwait" ;;
        "Asia/Qatar") 
            echo "25.2854,51.5310:timezone_qatar:Doha:Qatar" ;;
        "Asia/Bahrain") 
            echo "26.0667,50.5577:timezone_bahrain:Manama:Bahrain" ;;
        
        # South America (excluding Brazil)
        "America/Argentina/Buenos_Aires") 
            echo "-34.6118,-58.3960:timezone_argentina:Buenos Aires:Argentina" ;;
        "America/Lima") 
            echo "-12.0464,-77.0428:timezone_peru:Lima:Peru" ;;
        "America/Bogota") 
            echo "4.7110,-74.0721:timezone_colombia:Bogotá:Colombia" ;;
        "America/Caracas") 
            echo "10.4806,-66.9036:timezone_venezuela:Caracas:Venezuela" ;;
        "America/Santiago") 
            echo "-33.4489,-70.6693:timezone_chile:Santiago:Chile" ;;
        
        # Default fallback
        *) 
            echo "0.0000,0.0000:timezone_unknown:Unknown:Unknown" ;;
    esac
    
    log WARN "Using timezone-based approximation for $timezone"
}

# Generate location-based node ID with better encoding
generate_location_id() {
    local coords_with_method="$1"
    local coords=$(echo "$coords_with_method" | cut -d':' -f1)
    local method=$(echo "$coords_with_method" | cut -d':' -f2)
    local city=$(echo "$coords_with_method" | cut -d':' -f3)
    
    log DEBUG "Generating location ID for: $coords_with_method"
    
    # Create readable location ID
    local lat=$(echo "$coords" | cut -d',' -f1 | tr -d '-.' | cut -c1-4)
    local lon=$(echo "$coords" | cut -d',' -f2 | tr -d '-.' | cut -c1-4)
    
    # Clean city name for ID - remove special characters and limit length
    local city_clean=$(echo "$city" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]//g' | cut -c1-8)
    
    # Method prefix
    local method_prefix=""
    case "$method" in
        "ip_geolocation"*) method_prefix="geo" ;;
        "timezone"*) method_prefix="tz" ;;
        "manual") method_prefix="man" ;;
        *) method_prefix="unk" ;;
    esac
    
    # Generate short random suffix
    local random_suffix=$(openssl rand -hex 3)
    
    # Ensure valid city name, fallback to coordinates
    if [ -z "$city_clean" ] || [ "$city_clean" = "unknown" ]; then
        city_clean="loc${lat}${lon}"
    fi
    
    local location_id="${method_prefix}-${city_clean}-${random_suffix}"
    log DEBUG "Generated location ID: $location_id"
    
    echo "$location_id"
}

# Validate coordinates format
validate_coordinates() {
    local coords="$1"
    
    if [[ ! "$coords" =~ ^-?[0-9]+\.?[0-9]*,-?[0-9]+\.?[0-9]*$ ]]; then
        log ERROR "Invalid coordinate format: $coords"
        echo "Expected format: latitude,longitude (e.g., -23.5505,-46.6333)"
        return 1
    fi
    
    local lat=$(echo "$coords" | cut -d',' -f1)
    local lon=$(echo "$coords" | cut -d',' -f2)
    
    # Validate latitude range (-90 to 90)
    if (( $(echo "$lat < -90 || $lat > 90" | bc -l) )); then
        log ERROR "Invalid latitude: $lat (must be between -90 and 90)"
        return 1
    fi
    
    # Validate longitude range (-180 to 180)
    if (( $(echo "$lon < -180 || $lon > 180" | bc -l) )); then
        log ERROR "Invalid longitude: $lon (must be between -180 and 180)"
        return 1
    fi
    
    log DEBUG "Coordinates validated: $coords"
    return 0
}

# Get geographical information summary
get_location_summary() {
    local coords_with_info="$1"
    local coords=$(echo "$coords_with_info" | cut -d':' -f1)
    local method=$(echo "$coords_with_info" | cut -d':' -f2)
    local city=$(echo "$coords_with_info" | cut -d':' -f3)
    local country=$(echo "$coords_with_info" | cut -d':' -f4)
    
    echo "Geographic Information:"
    echo "  Coordinates: $coords"
    echo "  City: $city"
    echo "  Country: $country"
    echo "  Detection Method: $method"
    
    case "$method" in
        "ip_geolocation"*)
            echo "  Accuracy: High (IP-based geolocation)"
            ;;
        "timezone"*)
            echo "  Accuracy: Medium (timezone approximation)"
            ;;
        "manual")
            echo "  Accuracy: Exact (manually specified)"
            ;;
        *)
            echo "  Accuracy: Unknown"
            ;;
    esac
}