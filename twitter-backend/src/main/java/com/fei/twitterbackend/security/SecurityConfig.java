package com.fei.twitterbackend.security;

import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;

import java.util.List;

@Configuration
@EnableWebSecurity
@RequiredArgsConstructor
public class SecurityConfig {
    @Value("${app.frontend.url}")
    private String frontendUrl;

    private final JwtAuthenticationFilter jwtAuthFilter;

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
                // 1. Disable CSRF (Stateful protection not needed for stateless JWT APIs)
                .csrf(AbstractHttpConfigurer::disable)

                // 2. Configure CORS (Allow Next.js frontend)
                .cors(cors -> cors.configurationSource(corsConfigurationSource()))

                // 3. Define Endpoint Rules
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/error").permitAll()
                        // Allow Public Access to Auth Endpoints
                        .requestMatchers("/api/v1/auth/**").permitAll()
                        // All other requests require a valid JWT
                        .anyRequest().authenticated()
                )
                .exceptionHandling(e -> e
                        // 1. Handle "Not Logged In" (401)
                        .authenticationEntryPoint((request, response, authException) -> {
                            response.sendError(401, "Unauthorized: Please log in");
                        })
                        // 2. Handle "Logged In, But No Permission" (403) - Optional: Spring does this by default, but you can customize it here
                        .accessDeniedHandler((request, response, accessDeniedException) -> {
                            response.sendError(403, "Forbidden: You don't have permission");
                        })
                )

                // 4. Make it Stateless (No Session/Cookies created by Spring)
                .sessionManagement(sess -> sess.sessionCreationPolicy(SessionCreationPolicy.STATELESS))

                // 5. Add Custom JWT Filter BEFORE the standard Spring Login filter
                .addFilterBefore(jwtAuthFilter, UsernamePasswordAuthenticationFilter.class);

        return http.build();
    }

    // CORS Config
    @Bean
    public UrlBasedCorsConfigurationSource corsConfigurationSource() {
        CorsConfiguration config = new CorsConfiguration();
        config.setAllowCredentials(true);
        config.setAllowedOrigins(List.of(frontendUrl)); // allow frontend
        config.setAllowedHeaders(List.of("Authorization", "Content-Type"));
        config.setAllowedMethods(List.of("GET", "POST", "PUT", "DELETE", "OPTIONS"));

        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", config);
        return source;
    }
}