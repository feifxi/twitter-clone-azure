package com.fei.twitterbackend.security;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

@Configuration
@EnableWebSecurity
@RequiredArgsConstructor
public class SecurityConfig {

    private final JwtAuthenticationFilter jwtAuthFilter;

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
                // 1. Disable CSRF (Stateful protection not needed for stateless JWT APIs)
                .csrf(AbstractHttpConfigurer::disable)

                // 2. Define Endpoint Rules
                .authorizeHttpRequests(auth -> auth
                        // SPECIFIC RESTRICTIONS (The Exceptions)
                        .requestMatchers(HttpMethod.GET, "/api/v1/feeds/following").authenticated()

                        // Error and Auth are always public
                        .requestMatchers("/error").permitAll()
                        .requestMatchers("/api/v1/auth/**").permitAll()

                        // Public Read Access (GET)
                        // Allow guests to see the feed and single tweets
                        .requestMatchers(HttpMethod.GET, "/api/v1/feeds/**").permitAll()
                        .requestMatchers(HttpMethod.GET, "/api/v1/tweets/**").permitAll()
                        .requestMatchers(HttpMethod.GET, "/api/v1/users/**").permitAll()
                        .requestMatchers(HttpMethod.GET, "/api/v1/search/**").permitAll()
                        .requestMatchers(HttpMethod.GET, "/api/v1/discovery/**").permitAll()

                        // Everything else (POST, PUT, DELETE, and other paths)
                        // This covers creating tweets, liking, following, etc.
                        .anyRequest().authenticated()
                )

                // 3. Exception
                .exceptionHandling(e -> e
                        // Handle "Not Logged In" (401)
                        .authenticationEntryPoint((request, response, authException) -> {
                            response.sendError(401, "Unauthorized: Please log in");
                        })
                        // Handle "Logged In, But No Permission" (403) - Optional: Spring does this by default, but you can customize it here
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
}